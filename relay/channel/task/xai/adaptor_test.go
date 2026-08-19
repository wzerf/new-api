package xai

import (
	"net/http/httptest"
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEstimateBillingUsesDurationSeconds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Set("task_request", relaycommon.TaskSubmitReq{
		Duration: 10,
	})

	adaptor := &TaskAdaptor{}
	ratios := adaptor.EstimateBilling(ctx, nil)
	// default 720p multiplier is 0.7
	require.Contains(t, ratios, "total_units")
	assert.InDelta(t, 7.0, ratios["total_units"], 1e-9)
}

func TestEstimateBillingMetadataOverridesRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Set("task_request", relaycommon.TaskSubmitReq{
		Duration: 6,
		Metadata: map[string]any{
			"duration": 12,
			"resolution": "1080p",
		},
	})

	adaptor := &TaskAdaptor{}
	ratios := adaptor.EstimateBilling(ctx, nil)
	require.Contains(t, ratios, "total_units")
	assert.InDelta(t, 12.0, ratios["total_units"], 1e-9)
}

func TestResolveDurationSecondsCapsMetadata(t *testing.T) {
	got := ResolveDurationSeconds(map[string]any{"duration": relaycommon.MaxTaskDurationSeconds + 100}, 6, "4")
	assert.Equal(t, relaycommon.MaxTaskDurationSeconds, got)
}

func TestParseTaskResultSuccessUsesVideoURL(t *testing.T) {
	adaptor := &TaskAdaptor{}
	ti, err := adaptor.ParseTaskResult([]byte(`{"status":"done","video":{"url":"https://cdn.example/video.mp4","duration":8}}`))
	require.NoError(t, err)
	assert.Equal(t, "https://cdn.example/video.mp4", ti.Url)
	assert.Equal(t, "100%", ti.Progress)
}
