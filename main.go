package main

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var Interrupted = false

func UseContext(context string) error {
	return Run("kubectl", []string{"config", "use-context", context})
}
func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func ExecForEach(context, namespace, kubeconfig string, args []string) {
	lg := log.WithFields(log.Fields{
		"context":   context,
		"namespace": namespace,
	})
	if Interrupted {
		lg.Info("Interrupted closing")
		return
	}
	lg.Info("===============================================")
	if err := UseContext(context); err != nil {
		lg.WithError(err).Error("failed changing context")
	}
	if err := Run("kubectl", args); err != nil {
		lg.WithError(err).Error("failed executing cmd")
	}
}

// deleteEmptyFields remove empty string from slice
func deleteEmptyFields(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// Run will execute commands
func Run(command string, args []string) error {

	args = deleteEmptyFields(args)

	cmd := exec.Command(command, args...)
	var stderr bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)

	if err != nil && elapsed < time.Second {
		errStr := stderr.String()
		log.WithFields(log.Fields{
			"command": command,
			"args":    args,
		}).Error(errStr)
	}

	return err

}

func PrintHelp() {
	log.Info("Usage: kq get pods --namespace kube-system")
}

func ValidateAndGet() ([]string, error) {
	if len(os.Args) < 2 {
		return nil, errors.New("Minimum 2 params required")
	}
	result := os.Args[1:]
	if result[0] == "help" || result[0] == "h" || result[0] == "-h" || result[0] == "--help" {
		return nil, errors.New("")
	}
	return result, nil
}

func BackToOriginalNamespace(currCtx, currNs string) {
	Interrupted = true
	Run("kubectl", []string{"config", "use-context", currCtx, "--namespace", currNs})
}
func SetupCloseHandler(currCtx, currNs string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("\r- Shutting down, changing back to the original context:", currCtx)
		BackToOriginalNamespace(currCtx, currNs)
		os.Exit(0)
	}()
}

func main() {

	params, err := ValidateAndGet()
	if err != nil {
		log.WithError(err).Error(err)
		PrintHelp()
		return
	}

	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	}

	flag.Parse()

	clientCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	currCtx := clientCfg.CurrentContext
	currNs := clientCfg.Contexts[currCtx].Namespace

	SetupCloseHandler(currCtx, currNs)
	for _, ctx := range clientCfg.Contexts {
		if strings.Contains(ctx.Cluster, "arn:aws:eks") {
			ExecForEach(ctx.Cluster, ctx.Namespace, *kubeconfig, params)
		}
	}
	BackToOriginalNamespace(currCtx, currNs)
}
