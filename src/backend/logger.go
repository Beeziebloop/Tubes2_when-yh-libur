package main
import(
	"fmt"
	"os"
	"strings"
	"time"
)

func PrintTraversalLogResult(result TraversalRes, selector string, filename string) error{
	f, err := os.Create(filename)
	if err != nil{
		return fmt.Errorf("Gagal bikin logfile: %w", err)
	}
	defer f.Close()

	algo := "unknown"
	if len(result.TraversalLog) > 0{
		algo = result.TraversalLog[0].Algorithm
	}

	//header
	fmt.Fprintln(f, "==========Traversal Log==========")
	fmt.Fprintf(f, "Algorithm: %s\n", algo)
	fmt.Fprintf(f, "Selector: %s\n", selector)
	fmt.Fprintf(f, "Timestamp: %s\n", time.Now().Format("2001-09-11 10:00:67"))
	fmt.Fprintln(f, strings.Repeat("-", 50))
	//steps
	for _, entry := range result.TraversalLog{
		classes := strings.Join(entry.N_Classes, " ")
		if classes == ""{
			classes = "-"
		}
		matchedString := "no"
		if entry.Is_matched{
			matchedString = "yes"
		}
		fmt.Fprintf(f, "Step %03d | Tag: %-8s | ID: %-12s | Classes: %-20s | Depth: %d | Matched: %s\n", entry.Step, entry.N_Tag, ifEmpty(entry.N_ID), classes, entry.Depth, matchedString)
	}
	//summarynya
	fmt.Fprintln(f, strings.Repeat("-", 50))
	fmt.Fprintln(f, "==========Summary==========")
	fmt.Fprintf(f, "Total visited: %d\n", result.VisCount)
	fmt.Fprintf(f, "Total matched: %d\n", len(result.MatchedNodes))
	fmt.Fprintf(f, "Elapsed time: %s\n", result.ElapsedTime)
	fmt.Fprintf(f, "Max depth visited in process: %d\n", result.MaxDepthVisited)
	fmt.Fprintf(f, "Max depth of the whole tree: %d\n", result.FullMaxDepth)
	return nil;
}

func ifEmpty(s string) string{
	if s == ""{
		return "-"
	}
	return s
}