package catalog

import (
	"context"
	"testing"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
)

func TestCatalogTorrentViewPolicy(t *testing.T) {
	t.Parallel()

	domain := NewCatalogTorrentDomain()
	owner := &model.Actor{Id: 10}
	member := &model.Actor{Id: 20}
	staff := &model.Actor{Id: 30, IsStaff: true}

	tests := []struct {
		name    string
		actor   *model.Actor
		status  int
		banned  bool
		allowed bool
	}{
		{name: "anonymous can view published", status: consts.CatalogTorrentStatusPublished, allowed: true},
		{name: "member can view published", actor: member, status: consts.CatalogTorrentStatusPublished, allowed: true},
		{name: "owner can view pending", actor: owner, status: consts.CatalogTorrentStatusPending, allowed: true},
		{name: "owner can view rejected", actor: owner, status: consts.CatalogTorrentStatusRejected, allowed: true},
		{name: "member cannot view pending", actor: member, status: consts.CatalogTorrentStatusPending},
		{name: "member cannot view rejected", actor: member, status: consts.CatalogTorrentStatusRejected},
		{name: "staff can view pending", actor: staff, status: consts.CatalogTorrentStatusPending, allowed: true},
		{name: "staff can view rejected", actor: staff, status: consts.CatalogTorrentStatusRejected, allowed: true},
		{name: "owner cannot view banned", actor: owner, status: consts.CatalogTorrentStatusPublished, banned: true},
		{name: "staff can inspect banned", actor: staff, status: consts.CatalogTorrentStatusPublished, banned: true, allowed: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			torrent := &entity.CatalogTorrent{OwnerId: owner.Id, Status: tt.status, Banned: tt.banned}
			err := domain.CheckTorrentViewPolicy(context.Background(), tt.actor, torrent)
			if (err == nil) != tt.allowed {
				t.Fatalf("allowed = %v, error = %v", err == nil, err)
			}
		})
	}
}

func TestCatalogTorrentDownloadPolicy(t *testing.T) {
	t.Parallel()

	domain := NewCatalogTorrentDomain()
	owner := &model.Actor{Id: 10}
	member := &model.Actor{Id: 20}
	staff := &model.Actor{Id: 30, IsStaff: true}

	tests := []struct {
		name    string
		actor   *model.Actor
		status  int
		banned  bool
		allowed bool
	}{
		{name: "anonymous cannot download", status: consts.CatalogTorrentStatusPublished},
		{name: "member can download published", actor: member, status: consts.CatalogTorrentStatusPublished, allowed: true},
		{name: "owner can download pending", actor: owner, status: consts.CatalogTorrentStatusPending, allowed: true},
		{name: "owner can download rejected", actor: owner, status: consts.CatalogTorrentStatusRejected, allowed: true},
		{name: "member cannot download pending", actor: member, status: consts.CatalogTorrentStatusPending},
		{name: "staff can download pending", actor: staff, status: consts.CatalogTorrentStatusPending, allowed: true},
		{name: "banned blocks owner", actor: owner, status: consts.CatalogTorrentStatusPublished, banned: true},
		{name: "banned blocks staff", actor: staff, status: consts.CatalogTorrentStatusPublished, banned: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			torrent := &entity.CatalogTorrent{OwnerId: owner.Id, Status: tt.status, Banned: tt.banned}
			err := domain.CheckTorrentDownloadPolicy(context.Background(), tt.actor, torrent)
			if (err == nil) != tt.allowed {
				t.Fatalf("allowed = %v, error = %v", err == nil, err)
			}
		})
	}
}

func TestCatalogTorrentAnnouncePolicy(t *testing.T) {
	t.Parallel()

	domain := NewCatalogTorrentDomain()
	owner := &model.Actor{Id: 10}
	member := &model.Actor{Id: 20}
	staff := &model.Actor{Id: 30, IsStaff: true}

	tests := []struct {
		name    string
		actor   *model.Actor
		status  int
		banned  bool
		allowed bool
	}{
		{name: "anonymous cannot announce", status: consts.CatalogTorrentStatusPublished},
		{name: "member can announce published", actor: member, status: consts.CatalogTorrentStatusPublished, allowed: true},
		{name: "owner can announce pending", actor: owner, status: consts.CatalogTorrentStatusPending, allowed: true},
		{name: "member cannot announce pending", actor: member, status: consts.CatalogTorrentStatusPending},
		{name: "staff can announce pending", actor: staff, status: consts.CatalogTorrentStatusPending, allowed: true},
		{name: "owner cannot announce rejected", actor: owner, status: consts.CatalogTorrentStatusRejected},
		{name: "staff cannot announce rejected", actor: staff, status: consts.CatalogTorrentStatusRejected},
		{name: "banned blocks owner", actor: owner, status: consts.CatalogTorrentStatusPublished, banned: true},
		{name: "banned blocks staff", actor: staff, status: consts.CatalogTorrentStatusPublished, banned: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			torrent := &entity.CatalogTorrent{OwnerId: owner.Id, Status: tt.status, Banned: tt.banned}
			err := domain.CheckTorrentAnnouncePolicy(context.Background(), tt.actor, torrent)
			if (err == nil) != tt.allowed {
				t.Fatalf("allowed = %v, error = %v", err == nil, err)
			}
		})
	}
}
