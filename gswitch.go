package main

import (
	"encoding/json"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os/exec"
	"os"
	"strings"

	. "github.com/logrusorgru/aurora"
)

type AccountFile struct {
	ProjectId string `json:"project_id"`
}

type Cluster struct {
	Name string `json:"name"`
	Zone string `json:"zone"`
}

func empty(s string) bool {
    return len(strings.TrimSpace(s)) == 0
}

func isInstalled(tool string) bool {
	cmd := exec.Command(tool, "version")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func gcloudAuth(noLaunchBrowser bool) {
	var noBrowser string = ""
	if noLaunchBrowser {
		noBrowser = "--no-launch-browser"
	}
	cmd := exec.Command("gcloud", "auth", "login", noBrowser)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("gcloud authentication request finished with error: %v", err)
	}
}

func gcloudSetProject(projectId string) {
	cmd := exec.Command("gcloud", "config",	"set", "project", projectId)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("gcloud set project finished with error: %v", err)
	}
}

func gcloudSetConfiguration(configName string) {
	cmd := exec.Command(
		"gcloud",
		"config",
		"configurations",
		"activate",
		configName,
	)
	log.Println("Activating " + configName + " gcloud configuration...")
	err := cmd.Run()
	if err != nil {
		printWarning(Sprintf("Can't switch to %s configuration. Command finished with error: %v", configName, err))
	}
}

func gcloudConfigure(accountFile, projectId string) {
	cmd := exec.Command(
		"gcloud", 
		"auth", 
		"activate-service-account", 
		"--key-file", accountFile, 
		"--project", projectId,
	)
	log.Println("Modifying gcloud configuration...")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("gcloud configuration finished with error: %v", err)
	}
}

func getKubernetesClusters() []Cluster {
	out, err := exec.Command(
		"gcloud", 
		"container", 
		"clusters", 
		"list",
		"--format=json",
	).Output()
	if err != nil {
		raiseError(Sprintf("'gcloud container clusters list' finished with error: %v", err))
	}
	var clusters []Cluster
	json.Unmarshal(out, &clusters)
	
	return clusters
}

func kubectlConfigure(clusterName, clusterZone string) {
	cmd := exec.Command(
		"gcloud", 
		"container", 
		"clusters", 
		"get-credentials", clusterName, 
		"--zone", clusterZone,
	)
	log.Println("Modifying kubectl configuration...")
	err := cmd.Run()
	if err != nil {
		raiseError(Sprintf("kubectl configuration finished with error: %v", err))
	}
}

func raiseError(message string) {
	log.Fatal(Sprintf("%s %s", Bold(Red("Error:")), Bold(message)))
}

func printWarning(message string) {
	log.Println(Sprintf("%s %s", Bold(Yellow("Warning:")), Bold(message)))
}

func serviceAccountFlow() string {
	accountPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if empty(accountPath) {
		errorMessage := "Could not find default credentials. " +
						"See https://developers.google.com/accounts/docs/application-default-credentials " +
						"for more information"
		raiseError(errorMessage)
	}

    credentials, err := os.Open(accountPath)
    if err != nil {
        raiseError(Sprintf(err))
	}
	defer credentials.Close()

	byteValue, _ := ioutil.ReadAll(credentials)

	var account AccountFile
	json.Unmarshal([]byte(byteValue), &account)

	gcloudConfigure(accountPath, account.ProjectId)

	return account.ProjectId
}

func main() {
	var projectId string
	flag.StringVarP(&projectId, "project", "p", "", "Project ID for gcloud configuration. " + 
					"If project ID set, you can't use service account auth")
	usePreviousIdentity := flag.Bool("no-auth", false, "Do not authenticate user. Use previous identity")
	useServiceAccount := flag.Bool("use-service-account", false, "Use service account instead of Google account. " + 
								   "This option overwrites any other options!")
	noLaunchBrowser := flag.Bool("no-launch-browser", false, "Do not launch a browser for authorization. " +
						   		 "If enabled or DISPLAY variable is not set, prints a URL to standard output to be copied. " +
						   		 "Disabled by default")
	flag.Parse()
		
	if !isInstalled("gcloud") {
		raiseError("gcloud is not installed or can't be found in PATH. Nothing to configure")
	}

	gcloudSetConfiguration("default")
	
	if *useServiceAccount {
		projectId = serviceAccountFlow()
	} else {
		if !*usePreviousIdentity {
			gcloudAuth(*noLaunchBrowser)
		}
		if projectId != "" {
			gcloudSetProject(projectId)
		}
	}

	if !isInstalled("kubectl") {
		printWarning("kubectl is not installed. Skipping its configuration")
		os.Exit(0)
	}

	clusters := getKubernetesClusters()

	numberOfClusters := len(clusters)

	if numberOfClusters > 1 {
		printWarning("More than one cluster is available to configure. Configuring first one")
	} else if numberOfClusters == 0 {
		printWarning("No kubernetes clusters discovered in project " + projectId + ". Skipping kubectl configuration")
		os.Exit(0)
	}

	cluster := clusters[0]

	kubectlConfigure(cluster.Name, cluster.Zone)

	log.Println(Sprintf("%s %s", Bold(Green("gTools've been switched to")), Bold(Green(projectId))))
}
