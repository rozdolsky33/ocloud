#!/bin/bash
# trigger
# Update all references to pkg/flags to internal/config
find . -type f -name "*.go" -exec sed -i '' 's/github.com\/rozdolsky33\/ocloud\/pkg\/flags/github.com\/rozdolsky33\/ocloud\/internal\/config/g' {} \;
find . -type f -name "*.go" -exec sed -i '' 's/flags\./config\./g' {} \;