// +build bin

package cmd

var (
	// UpgradeLater is the flag value.
	UpgradeLater bool
	// DonotUpgrade is the flag value.
	DonotUpgrade bool
)

func init() {
	OstentCmd.PersistentFlags().BoolVar(&UpgradeLater, "upgradelater", false, "Delay startup upgrade check and applying available upgrade")
	OstentCmd.PersistentFlags().BoolVar(&DonotUpgrade, "noupgrade", false, "Do not upgrade, but log if there's an upgrade")
}
