package main

import (
	"errors"
	"github.com/bndr/gojenkins"
	"github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type JenkinsServerConfig struct {
	Server string `yaml:"server"`
	User   string `yaml:"user"`
	Token  string `yaml:"token"`
}

type JenkinsJobConfig struct {
	Name    string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	Options map[string]string `yaml:"options"`
}

type Config struct {
	Jenkins JenkinsServerConfig `yaml:"jenkins"`
	Jobs []JenkinsJobConfig `yaml:"jobs"`
}

func triggerJenkinsJob(serverConfig JenkinsServerConfig, jenkinsToken string,  jobConfig JenkinsJobConfig ) error {

	jenkins := gojenkins.CreateJenkins(nil, serverConfig.Server, serverConfig.User, jenkinsToken)
	_, err := jenkins.Init()
	if err != nil {
		return err
	}
	_, err = jenkins.BuildJob(jobConfig.Name, jobConfig.Options)
	if err != nil {
		return err
	}
	return nil
}

func getJenkinsToken(serverConfig JenkinsServerConfig) (string, error) {
	jenkinsToken := os.Getenv("JENKINS_TOKEN")
	if jenkinsToken == "" {
		jenkinsToken = serverConfig.Token
	}

	if jenkinsToken != "" {
		return jenkinsToken, nil
	}

	return "", errors.New("unable to get Jenkins token")
}

func main()  {
	configFileByte, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		logrus.Fatal("unable to read config file", err)
	}
	config := &Config{}

	err = yaml.Unmarshal(configFileByte, config)
	if err != nil {
		logrus.Fatal("unable to unmarshal config data ", err)
	}

	jenkinsToken, err := getJenkinsToken(config.Jenkins)
	if err != nil {
		logrus.Fatalln(err)
	}

	c := cron.New()
	for _, job := range config.Jobs {
		jobConfig := job
		_, err = c.AddFunc(job.Schedule, func() {
			err := triggerJenkinsJob(config.Jenkins, jenkinsToken, jobConfig )
			if err != nil {
				logrus.Fatalln("unable to trigger jenknis job", err)
			}
		})

		if err != nil {
			logrus.Fatalln("Error adding new job", err)

		}
	}

	c.Start()

	time.Sleep(time.Second * 100000000)
}