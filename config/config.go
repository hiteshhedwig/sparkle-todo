package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
)

const localfile = "./config.cfg"
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type CompletedTask struct {
	Index int
	Task  string
	Time  string //time at which it got complete
}

type CurrentTasks struct {
	Index int
	Task  string
	Time  string
}

type CancelledTask struct {
	Index int
	Task  string
	Time  string //time at which it got cancelled
}

type Config struct {
	Project        string
	Currenttasks   []CurrentTasks
	Completedtasks []CompletedTask
	Cancelledtasks []CancelledTask
}

type Conversion interface {
	ToCancelledTask() *CancelledTask
	ToCurrentTask() *CurrentTasks
	ToCompletedTask() *CompletedTask
}

func (c *Config) SetProjectName(name string) error {
	fmt.Printf("Setting Project name from : %s ", c.Project)
	c.Project = name
	fmt.Printf("-> %s", name)
	c.SaveConfig()
	return nil
}

func (c *Config) CompletedTask(idx int) error {
	// idx check!
	if idx > len(c.Currenttasks) {
		fmt.Printf("%s Index value [%d] is not correct, Please check again in below table \n \n", string(colorRed), idx)
		fmt.Printf("")
		c.ListAllCurrentTasks()
		return fmt.Errorf("Out of range")
	}

	donetask := c.Currenttasks[idx-1]
	disp := emoji.Sprintf(":pizza:")
	fmt.Printf(fmt.Sprintf("%s", disp))
	fmt.Printf("%s You have completed the task : %s \n", string(colorGreen), donetask.Task)

	c.Currenttasks = remove(c.Currenttasks, idx-1)
	can := c.ToCompletedTask(donetask)
	c.Completedtasks = append(c.Completedtasks, *can)

	for i := range c.Currenttasks {
		c.Currenttasks[i].Index = i + 1
	}

	c.SaveConfig()
	return nil
}

func (c *Config) ToCompletedTask(com CurrentTasks) *CompletedTask {

	return &CompletedTask{
		Index: com.Index,
		Task:  com.Task,
		Time:  time.Now().Format("2006-01-02 3:4:5 pm"),
	}

}

func (c *Config) ToCancelledTask(com CurrentTasks) *CancelledTask {

	return &CancelledTask{
		Index: com.Index,
		Task:  com.Task,
		Time:  time.Now().Format("2006-01-02 3:4:5 pm"),
	}

}

func (c *Config) Cancel(idx int) {

	var userin string

	dumptask := c.Currenttasks[idx-1]
	log.Info(dumptask)
	fmt.Println("Are you sure you want to cancel this task? ", dumptask.Task)
	fmt.Print("Type yes/no to confirm cancellation : ")
	fmt.Scanf("%s", &userin)

	switch strings.ToLower(userin) {
	case "no":
		fmt.Println("Aborting")
		return
	case "yes":
		fmt.Println("ok to cancel this task")
		c.CancelforReal(idx)
	default:
		fmt.Println("Aborting")
		return
	}

}

func (c *Config) CancelforReal(idx int) {

	dumptask := c.Currenttasks[idx-1]
	c.Currenttasks = remove(c.Currenttasks, idx-1)
	fmt.Println(dumptask)
	can := c.ToCancelledTask(dumptask)
	c.Cancelledtasks = append(c.Cancelledtasks, *can)

	for i := range c.Currenttasks {
		c.Currenttasks[i].Index = i + 1
	}

	c.SaveConfig()
}

func remove(slice []CurrentTasks, s int) []CurrentTasks {
	return append(slice[:s], slice[s+1:]...)
}

func (c *Config) ListAllCurrentTasks() {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Index", "Task", "Time"})

	for _, currenttask := range c.Currenttasks {
		t.AppendSeparator()
		t.AppendRow([]interface{}{currenttask.Index, currenttask.Task, currenttask.Time})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()

}

func (c *Config) ListTasksCatWise() {

	var typetasks [][]interface{}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Index", "Task", "Time"})

	//typetask := []interface{}{000, "Current Task", " --"}
	typetasks = append(typetasks, []interface{}{"--", "Current Task", " --"})
	typetasks = append(typetasks, []interface{}{"--", "Completed Task", " --"})
	typetasks = append(typetasks, []interface{}{"--", "Cancelled Task", " --"})

	for idx, typetask := range typetasks {
		t.AppendSeparator()
		t.AppendRow(typetask)
		t.AppendSeparator()

		switch idx {
		case 0:
			for _, tsk := range c.Currenttasks {
				t.AppendRow([]interface{}{tsk.Index, tsk.Task, tsk.Time})
			}
		case 1:
			for _, tsk := range c.Completedtasks {
				t.AppendRow([]interface{}{tsk.Index, tsk.Task, tsk.Time})
			}
		case 2:
			for _, tsk := range c.Cancelledtasks {
				t.AppendRow([]interface{}{tsk.Index, tsk.Task, tsk.Time})
			}
		}

	}
	t.SetStyle(table.StyleRounded)
	t.Render()

}

func (c *Config) AddToCurrentTasks(task string) {

	var index int

	switch len(c.Currenttasks) {
	case 0:
		index = 1
	default:
		index = len(c.Currenttasks) + 1
	}

	curr := CurrentTasks{
		Task:  task,
		Index: index,
		Time:  time.Now().Format("2006-01-02 3:4:5 pm"),
	}
	c.Currenttasks = append(c.Currenttasks, curr)
}

func (config Config) SaveConfig() {

	configBytes, err := json.MarshalIndent(config, "", "   ")
	if err != nil {
		fmt.Print("Something Wrong happened", err)
	}

	err = ioutil.WriteFile(localfile, configBytes, os.ModePerm)

	if err != nil {
		fmt.Print("Error writing device config!", err)
	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func LoadConfig(project string) (*Config, error) {

	var config Config

	if !fileExists(localfile) {
		fmt.Printf("Setting up config file! \n")
		config.Project = project
		configBytes, err := json.MarshalIndent(config, "", "   ")
		if err != nil {
			fmt.Print("Something Wrong happened", err)
			return nil, err
		}

		err = ioutil.WriteFile(localfile, configBytes, os.ModePerm)

		if err != nil {
			fmt.Print("Error writing device config!", err)
			return nil, err
		}
	} else {
		configFile, err := os.Open(localfile)
		if err != nil {
			log.Error("Error opening device file", err)
			return nil, err
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer configFile.Close()
		configBytes, _ := ioutil.ReadAll(configFile)
		err = json.Unmarshal(configBytes, &config)

		if err != nil {
			log.Error("11Error loading  config!", err)
			return nil, err
		}
	}

	return &config, nil
}
