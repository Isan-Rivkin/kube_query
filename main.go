package main

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

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

	lg.Info("================================================")
	//conf, err := buildConfigFromFlags(context, kubeconfig)
	//client, err := kubernetes.NewForConfig(conf)
	// if err != nil {
	// 	lg.WithError(err).Error("error getting kubeconfig")
	// 	panic(err.Error())
	// }
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
	// log.WithFields(log.Fields{
	// 	"command": strings.Join(append([]string{command}, args...), " "),
	// }).Info("execute command")

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
	log.Info("Usage: kq get pods -n pe")
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

func main() {
	//fmt.Println(len(os.Args), os.Args)
	params, err := ValidateAndGet()
	if err != nil {
		log.WithError(err).Error(err)
		PrintHelp()
		return
	}

	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// config, err = buildConfigFromFlags("", *kubeconfig)
	// fmt.Println(config)
	// if err != nil {
	// 	panic(err)
	// }
	// create the client
	// client, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err.Error())
	// }
	clientCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	currCtx := clientCfg.CurrentContext
	for _, ctx := range clientCfg.Contexts {
		if strings.Contains(ctx.Cluster, "arn:aws:eks") {
			ExecForEach(ctx.Cluster, ctx.Namespace, *kubeconfig, params)
		}
	}

	Run("kubectl", []string{"config", "use-context", currCtx})
}
