package contextops

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_MergeForDone(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("It closes the context immediately if no contexts are provided", func(t *testing.T) {
		ctx := MergeForDone()
		g.Eventually(ctx.Done).Should(BeClosed())
	})

	t.Run("It closes the merged context when any channels are closed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		oneCtx := MergeForDone(context.Background(), context.TODO(), ctx)
		g.Consistently(oneCtx.Done).ShouldNot(BeClosed())

		cancel()
		g.Eventually(oneCtx.Done).Should(BeClosed())
	})
}
