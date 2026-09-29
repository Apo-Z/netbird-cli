package main

import (
	"github.com/spf13/cobra"
)

var jobsGetCmd = &cobra.Command{
	Use:               "jobs <peer-name|id>",
	Aliases:           []string{"jb"},
	Short:             "List peer jobs",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: validArgsFunc(peerNames),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolvePeerID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		jobs, err := c.GetPeerJobs(id)
		if err != nil {
			printErr(err)
			return
		}
		printOutput(jobs)
	},
}

var jobGetCmd = &cobra.Command{
	Use:   "job <peer-name|id> <job-id>",
	Short: "Show a job",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := c.ResolvePeerID(args[0])
		if err != nil {
			printErr(err)
			return
		}
		job, err := c.GetPeerJob(id, args[1])
		if err != nil {
			printErr(err)
			return
		}
		printOutput(job)
	},
}

func init() {
	getCmd.AddCommand(jobsGetCmd, jobGetCmd)
}
