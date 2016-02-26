package integration_test

import (
	"fmt"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Fly CLI", func() {
	Describe("Pause Job", func() {
		var (
			flyCmd     *exec.Cmd
			reqsBefore int
		)

		pipelineName := "pipeline"
		jobName := "job-name-potato"
		fullJobName := fmt.Sprintf("%s/%s", pipelineName, jobName)

		BeforeEach(func() {
			flyCmd = exec.Command(flyPath, "-t", targetName, "pause-job", "-j", fullJobName)
			reqsBefore = len(atcServer.ReceivedRequests())
		})

		Context("when a job is paused using the API", func() {
			BeforeEach(func() {
				apiPath := fmt.Sprintf("/api/v1/pipelines/%s/jobs/%s/pause", pipelineName, jobName)
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", apiPath),
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)
			})

			It("successfully pauses the job", func() {
				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(0))
				Expect(atcServer.ReceivedRequests()).To(HaveLen(reqsBefore + 1))

				Eventually(sess).Should(gbytes.Say(fmt.Sprintf("paused '%s'", jobName)))
			})
		})

		Context("when a job is paused using the API", func() {
			BeforeEach(func() {
				apiPath := fmt.Sprintf("/api/v1/pipelines/%s/jobs/%s/pause", pipelineName, jobName)
				atcServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", apiPath),
						ghttp.RespondWith(http.StatusInternalServerError, nil),
					),
				)
			})

			It("exists 1 and outputs an error", func() {
				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess.Err).Should(gbytes.Say(`error`))

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(1))
				Expect(atcServer.ReceivedRequests()).To(HaveLen(reqsBefore + 1))
			})
		})
	})
})