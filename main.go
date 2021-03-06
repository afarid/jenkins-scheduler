package main

import (
	"errors"
	"github.com/bndr/gojenkins"
	log "github.com/sirupsen/logrus"
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
	Parameters map[string]string `yaml:"parameters"`
}

type Config struct {
	Jenkins JenkinsServerConfig `yaml:"jenkins"`
	Jobs []JenkinsJobConfig `yaml:"jobs"`
}

func triggerJenkinsJob(serverConfig JenkinsServerConfig, jobConfig JenkinsJobConfig ) error {

	jenkins := gojenkins.CreateJenkins(nil, serverConfig.Server, serverConfig.User, serverConfig.Token)
	_, err := jenkins.Init()
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"job-name": jobConfig.Name}).WithFields(log.Fields{"parameters": jobConfig.Parameters}).Info("Triggering Jenkins job")
	_, err = jenkins.BuildJob(jobConfig.Name, jobConfig.Parameters)
	if err != nil {
		return err
	}
	return nil
}

func setJenkinsToken(serverConfig *JenkinsServerConfig) error {
	jenkinsToken := os.Getenv("JENKINS_TOKEN")
	if jenkinsToken != "" {
		serverConfig.Token = jenkinsToken
	}
	if serverConfig.Token == "" {
		return errors.New("jenkins token cannot be empty, you need to add it to your config or as env variable JENKINS_TOKEN")
	}
	
	return  nil
}

func main()  {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Parsing config file")
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("unable to read config file", err)
	}
	config := &Config{}

	err = yaml.UnmarshalStrict(configFile, config)
	if err != nil {
		log.Fatal("unable to unmarshal config data ", err)
	}

	err = setJenkinsToken(&config.Jenkins)
	if err != nil {
		log.Fatalln(err)
	}

	c := cron.New()
	for _, job := range config.Jobs {
		jobConfig := job
		_, err = c.AddFunc(job.Schedule, func() {
			err := triggerJenkinsJob(config.Jenkins, jobConfig )
			if err != nil {
				log.Fatalln("unable to trigger jenknis job", err)
			}
		})

		if err != nil {
			log.Fatalln("Error adding new job", err)

		}
	}

	c.Start()

	time.Sleep(time.Second * 100000000)
}