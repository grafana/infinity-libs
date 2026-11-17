---
"@grafana/infinity-framesql": patch
---

🐛 **Bug Fix**: `min()` and `max()` no longer panic when a field has no numeric values (only nulls or text); they return null
