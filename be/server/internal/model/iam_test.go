package model

import (
	"math"
	"testing"

	"server/internal/model/entity"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func TestIamRankChainSkipsStaffAndLinksRoles(t *testing.T) {
	chain := NewIamRankChain([]entity.IamRole{
		{Id: 3, Level: 20},
		{Id: 1, Level: 0},
		{Id: 9, Level: 100, IsStaff: true},
		{Id: 2, Level: 10},
	})

	if chain.Len() != 3 {
		t.Fatalf("expected 3 normal rank nodes, got %d", chain.Len())
	}
	node := chain.NodeByRoleId(2)
	if node == nil {
		t.Fatal("expected role 2 node")
	}
	if node.Prev == nil || node.Prev.Role.Id != 1 {
		t.Fatalf("expected role 2 prev role 1, got %#v", node.Prev)
	}
	if node.Next == nil || node.Next.Role.Id != 3 {
		t.Fatalf("expected role 2 next role 3, got %#v", node.Next)
	}
	if staff := chain.NodeByRoleId(9); staff != nil {
		t.Fatalf("expected staff role to be skipped, got %#v", staff)
	}
}

func TestIamRankNodePromotionTargetUsesHighestMatchedRole(t *testing.T) {
	chain := NewIamRankChain([]entity.IamRole{
		rankRole(1, 0, nil),
		rankRole(2, 10, map[string]any{
			iamRankRulePromotion: []map[string]any{{iamRankFieldRatioGte: 1.0}},
		}),
		rankRole(3, 20, map[string]any{
			iamRankRulePromotion: []map[string]any{{iamRankFieldRatioGte: 2.0}},
		}),
	})

	result := chain.NodeByRoleId(1).ResolveTarget(IamRankStats{ratio: 2.5})
	if result.ShouldBan {
		t.Fatal("promotion should not request ban")
	}
	if result.Target == nil || result.Target.Role.Id != 3 {
		t.Fatalf("expected highest matched role 3, got %#v", result.Target)
	}
}

func TestIamRankNodeDemotionTargetUsesHighestLowerMatchedRole(t *testing.T) {
	chain := NewIamRankChain([]entity.IamRole{
		rankRole(1, 0, nil),
		rankRole(2, 10, map[string]any{
			iamRankRulePromotion: []map[string]any{{iamRankFieldRatioGte: 1.0}},
		}),
		rankRole(3, 20, map[string]any{
			iamRankRulePromotion: []map[string]any{{iamRankFieldRatioGte: 2.0}},
		}),
		rankRole(4, 30, map[string]any{
			iamRankRuleDemotion: []map[string]any{{iamRankFieldRatioLt: 3.0}},
		}),
	})

	result := chain.NodeByRoleId(4).ResolveTarget(IamRankStats{ratio: 1.5})
	if result.ShouldBan {
		t.Fatal("non-lowest demotion should not request ban")
	}
	if result.Target == nil || result.Target.Role.Id != 2 {
		t.Fatalf("expected highest lower matched role 2, got %#v", result.Target)
	}
}

func TestIamRankNodeDemotionFallsBackToLowestRole(t *testing.T) {
	chain := NewIamRankChain([]entity.IamRole{
		rankRole(1, 0, nil),
		rankRole(2, 10, map[string]any{
			iamRankRulePromotion: []map[string]any{{iamRankFieldRatioGte: 1.0}},
		}),
		rankRole(3, 20, map[string]any{
			iamRankRuleDemotion: []map[string]any{{iamRankFieldRatioLt: 2.0}},
		}),
	})

	result := chain.NodeByRoleId(3).ResolveTarget(IamRankStats{ratio: 0.5})
	if result.ShouldBan {
		t.Fatal("non-lowest demotion should not request ban")
	}
	if result.Target == nil || result.Target.Role.Id != 1 {
		t.Fatalf("expected fallback lowest role 1, got %#v", result.Target)
	}
}

func TestIamRankNodeLowestDemotionRequestsBan(t *testing.T) {
	chain := NewIamRankChain([]entity.IamRole{
		rankRole(1, 0, map[string]any{
			iamRankRuleDemotion: []map[string]any{{iamRankFieldRatioLt: 1.0}},
		}),
		rankRole(2, 10, nil),
	})

	result := chain.NodeByRoleId(1).ResolveTarget(IamRankStats{ratio: 0.5})
	if !result.ShouldBan {
		t.Fatal("expected lowest demotion to request ban")
	}
	if result.Target != nil {
		t.Fatalf("expected no target when banning, got %#v", result.Target)
	}
}

func TestNewIamRankStatsUsesInfiniteRatioWhenOnlyUploaded(t *testing.T) {
	stats := NewIamRankStats(IamRankCandidate{Uploaded: 1}, nil)
	if !math.IsInf(stats.ratio, 1) {
		t.Fatalf("expected positive infinite ratio, got %f", stats.ratio)
	}
}

func rankRole(id uint, level int, rules map[string]any) entity.IamRole {
	return entity.IamRole{
		Id:    id,
		Level: level,
		Rules: gjson.New(rules),
	}
}
