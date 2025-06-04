package root

import "embed"

//go:embed config/crd/* test/helper/observability/install/assets/*
var Assets embed.FS
