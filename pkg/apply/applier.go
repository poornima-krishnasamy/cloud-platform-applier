package apply

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Applier struct {
	RunResults chan<- Result
}

type Result struct {
	Response string
	Error    string
	Folder   string
}

const (
	nRoutines = 2
)

func (a *Applier) ApplyNamespaceDirs(wg *sync.WaitGroup, chunkFolder []string) {
	defer wg.Done()

	for _, folder := range chunkFolder {

		a.apply_kubernetes_files(folder)
		// a.apply_terraform(folder)
	}

}

func (a *Applier) apply_kubernetes_files(folder string) {

	// Location of kubeconfig file
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err.Error())
	// }

	dd, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	files, err := fileSystem.listFiles(folderpath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		f, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%q \n", string(f))

		decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(f), 100)
		for {
			var rawObj runtime.RawExtension
			if err = decoder.Decode(&rawObj); err != nil {
				break
			}

			obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
			unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
			if err != nil {
				log.Fatal(err)
			}

			unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

			gr, err := restmapper.GetAPIGroupResources(c.Discovery())
			if err != nil {
				log.Fatal(err)
			}

			mapper := restmapper.NewDiscoveryRESTMapper(gr)
			mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
			if err != nil {
				log.Fatal(err)
			}

			var dri dynamic.ResourceInterface
			if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
				if unstructuredObj.GetNamespace() == "" {
					unstructuredObj.SetNamespace("default")
				}
				dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
			} else {
				dri = dd.Resource(mapping.Resource)
			}

			if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
				log.Fatal(err)
			}
		}
		if err != io.EOF {
			log.Fatal("eof ", err)
		}
	}
	/* var outb, errb bytes.Buffer

	kubectlArgs := []string{"-n", filepath.Base(folder), "apply", "-f", folder}

	kubectlCommand := exec.Command("kubectl", kubectlArgs...)

	kubectlCommand.Stdout = &outb
	kubectlCommand.Stderr = &errb
	kubectlCommand.Run()
	*/

	// result := &Result{Response: errb.String(), Error: errb.String(), Folder: folder}

	// a.RunResults <- *result

}
func (a *Applier) apply_terraform(folder string) {

	// Get the value of an Environment Variable
	bucket := os.Getenv("PIPELINE_STATE_BUCKET")
	key_prefix := os.Getenv("PIPELINE_STATE_KEY_PREFIX")
	lock_table := os.Getenv("PIPELINE_TERRAFORM_STATE_LOCK_TABLE")
	region := os.Getenv("PIPELINE_STATE_REGION")
	cluster := os.Getenv("PIPELINE_CLUSTER")

	// //Checking that an environment variable is present or not.

	key := key_prefix + cluster + "/" + filepath.Base(folder) + "/terraform.tfstate"

	// // err := os.RemoveAll(filepath.Join(folder, ".terraform"))
	// // if err != nil {
	// // 	result = Result{Error: "Cant remove .terraform folders", Response: "", Folder: folder}
	// // 	results <- result
	// // 	return
	// // }

	var outb, errb bytes.Buffer

	kubectlArgs := []string{
		"init",
		fmt.Sprintf("%s=bucket=%s", "-backend-config", bucket),
		fmt.Sprintf("%s=key=%s", "-backend-config", key),
		fmt.Sprintf("%s=dynamodb_table=%s", "-backend-config", lock_table),
		fmt.Sprintf("%s=region=%s", "-backend-config", region)}

	Command := exec.Command("terraform", kubectlArgs...)

	Command.Dir = folder + "/resources"
	Command.Stdout = &outb
	Command.Stderr = &errb
	Command.Run()

	kubectlArgs = []string{"plan"}

	Command = exec.Command("terraform", kubectlArgs...)

	Command.Dir = folder + "/resources"
	Command.Stdout = &outb
	Command.Stderr = &errb
	Command.Run()

	kubectlArgs = []string{"apply"}

	Command = exec.Command("terraform", kubectlArgs...)

	Command.Dir = folder + "/resources"
	Command.Stdout = &outb
	Command.Stderr = &errb
	Command.Run()

	result := &Result{Response: outb.String(), Error: errb.String(), Folder: folder}

	a.RunResults <- *result
}
