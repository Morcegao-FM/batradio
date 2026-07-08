// Permite usar React.act nos testes de hooks.
;(globalThis as Record<string, unknown>).IS_REACT_ACT_ENVIRONMENT = true
