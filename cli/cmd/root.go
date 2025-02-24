package cmd

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
	pkgcmd "github.com/linkerd/linkerd2/pkg/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultLinkerdNamespace = "linkerd"

	smiExtensionName = "smi"
)

var (

	// special handling for Windows, on all other platforms these resolve to
	// os.Stdout and os.Stderr, thanks to https://github.com/mattn/go-colorable
	stdout = color.Output
	stderr = color.Error

	apiAddr               string // An empty value means "use the Kubernetes configuration"
	controlPlaneNamespace string
	kubeconfigPath        string
	kubeContext           string
	impersonate           string
	impersonateGroup      []string
	verbose               bool

	// These regexs are not as strict as they could be, but are a quick and dirty
	// sanity check against illegal characters.
	alphaNumDash = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
)

// NewCmdSMI returns a new SMI command
func NewCmdSMI() *cobra.Command {
	smiCmd := &cobra.Command{
		Use:   "smi",
		Short: "smi manages the SMI extension of Linkerd service mesh",
		Long:  `smi manages the SMI extension of Linkerd service mesh.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// enable / disable logging
			if verbose {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.PanicLevel)
			}

			if !alphaNumDash.MatchString(controlPlaneNamespace) {
				return fmt.Errorf("%s is not a valid namespace", controlPlaneNamespace)
			}

			return nil
		},
	}

	smiCmd.PersistentFlags().StringVarP(&controlPlaneNamespace, "linkerd-namespace", "L", defaultLinkerdNamespace, "Namespace in which Linkerd is installed")
	smiCmd.PersistentFlags().StringVar(&kubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests")
	smiCmd.PersistentFlags().StringVar(&kubeContext, "context", "", "Name of the kubeconfig context to use")
	smiCmd.PersistentFlags().StringVar(&impersonate, "as", "", "Username to impersonate for Kubernetes operations")
	smiCmd.PersistentFlags().StringArrayVar(&impersonateGroup, "as-group", []string{}, "Group to impersonate for Kubernetes operations")
	smiCmd.PersistentFlags().StringVar(&apiAddr, "api-addr", "", "Override kubeconfig and communicate directly with the control plane at host:port (mostly for testing)")
	smiCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Turn on debug logging")
	smiCmd.AddCommand(newCmdInstall())
	smiCmd.AddCommand(newCmdUninstall())
	smiCmd.AddCommand(newCmdVersion())
	smiCmd.AddCommand(newCmdCheck())

	// resource-aware completion flag configurations
	pkgcmd.ConfigureNamespaceFlagCompletion(
		smiCmd, []string{"linkerd-namespace"},
		kubeconfigPath, impersonate, impersonateGroup, kubeContext)

	pkgcmd.ConfigureKubeContextFlagCompletion(smiCmd, kubeconfigPath)
	return smiCmd
}
