package config

import (
	"flag"
)

var OperatorConfig Config

type Config struct {
	AtlasDomain          string
	EnableLeaderElection bool
	MetricsAddr          string
}

// ParseConfiguration fills the 'OperatorConfig' from the flags passed to the program
func ParseConfiguration() {
	OperatorConfig = Config{}
	flag.StringVar(&OperatorConfig.AtlasDomain, "atlas-domain", "https://cloud-qa.mongodb.com", "the Atlas URL domain name (no slash in the end).")
	flag.StringVar(&OperatorConfig.MetricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&OperatorConfig.EnableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	flag.Parse()
}
