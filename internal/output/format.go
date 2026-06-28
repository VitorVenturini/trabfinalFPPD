package output

import (
	"fmt"
	"io"

	"trabfinalfppd/internal/metrics"
)

func PrintResult(w io.Writer, result metrics.RunResult) {
	fmt.Fprintf(
		w,
		"mode=%s n=%d seed=%d nodes=%d processes=%d elapsed_sec=%.6f label=%q\n",
		result.Mode,
		result.N,
		result.Seed,
		result.Nodes,
		result.Processes,
		result.ElapsedSec,
		result.Label,
	)
	fmt.Fprintln(w, "verification:")
	fmt.Fprintf(w, "  c00=%.15f\n", result.Verification.TopLeft)
	fmt.Fprintf(w, "  c0n=%.15f\n", result.Verification.TopRight)
	fmt.Fprintf(w, "  cn0=%.15f\n", result.Verification.BottomLeft)
	fmt.Fprintf(w, "  cnn=%.15f\n", result.Verification.BottomRight)
	fmt.Fprintf(w, "  checksum=%.15f\n", result.Verification.Checksum)
}
