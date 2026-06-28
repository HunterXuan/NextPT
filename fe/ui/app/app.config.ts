export default defineAppConfig({
  ui: {
    colors: {
      primary: 'sky',
      neutral: 'slate'
    },
    button: {
      defaultVariants: {
        size: 'md'
      }
    },
    card: {
      slots: {
        root: 'rounded-lg'
      }
    },
    input: {
      defaultVariants: {
        size: 'md'
      }
    },
    textarea: {
      defaultVariants: {
        size: 'md'
      }
    }
  }
})
