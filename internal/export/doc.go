// Package export serialises a port snapshot into a portable format
// (JSON or CSV) that can be consumed by external tools, dashboards,
// or archival pipelines.
//
// # Usage
//
//	f, _ := os.Create("ports.json")
//	defer f.Close()
//
//	e, err := export.New(export.FormatJSON, f)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := e.Write(currentPortSet); err != nil {
//		log.Fatal(err)
//	}
//
// Supported formats are export.FormatJSON and export.FormatCSV.
package export
