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

func Test_MergeDone(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("It panics if the mainCtx is nil", func(t *testing.T) {
		g.Expect(func() { MergeDone(nil) }).To(Panic())
	})

	t.Run("It panics if any optional contexts are nil", func(t *testing.T) {
		g.Expect(func() { MergeDone(nil, context.Background(), nil) }).To(Panic())
	})

	t.Run("It allows just the mainCtx to be provided", func(t *testing.T) {
		ctx, cancel := MergeDone(context.Background())
		g.Expect(ctx).ToNot(BeNil())
		g.Expect(cancel).ToNot(BeNil())
	})

	t.Run("It closes the context when cancel is called on the mainCtx", func(t *testing.T) {
		ctx, cancel := MergeDone(context.Background())
		g.Expect(ctx).ToNot(BeNil())
		g.Expect(cancel).ToNot(BeNil())

		g.Consistently(ctx.Done()).ShouldNot(BeClosed())

		cancel()
		g.Eventually(ctx.Done()).Should(BeClosed())
	})

	t.Run("It closes the context when an optional context is canceld", func(t *testing.T) {
		optionalContext, optionalCancel := context.WithCancel(context.Background())

		ctx, cancel := MergeDone(context.Background(), optionalContext)
		g.Expect(ctx).ToNot(BeNil())
		g.Expect(cancel).ToNot(BeNil())

		g.Consistently(ctx.Done()).ShouldNot(BeClosed())

		optionalCancel()
		g.Eventually(ctx.Done()).Should(BeClosed())
	})

	t.Run("It preserves the Values on the mainCtx", func(t *testing.T) {
		valueCtx := context.WithValue(context.Background(), "test", "true") // fine for simple test

		ctx, cancel := MergeDone(valueCtx)
		defer cancel()

		g.Expect(ctx).ToNot(BeNil())
		g.Expect(cancel).ToNot(BeNil())
		g.Expect(ctx.Value("test")).To(Equal("true"))
	})
}
