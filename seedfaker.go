// Package seedfaker provides deterministic synthetic data generation via FFI.
// Same seed + same arguments = identical output across CLI, Node, Python, PHP, Ruby, Go.
package seedfaker

/*
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/bin/darwin-arm64
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/bin/darwin-x86_64
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/bin/linux-x86_64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/bin/linux-arm64
#cgo LDFLAGS: -lseedfaker_ffi
#include "seedfaker.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"runtime"
	"unsafe"
)

// SeedFaker generates deterministic synthetic data.
type SeedFaker struct {
	handle *C.SfFaker
}

// Options for creating a SeedFaker instance.
type Options struct {
	Seed     string `json:"seed,omitempty"`
	Locale   string `json:"locale,omitempty"`
	Tz       string `json:"tz,omitempty"`
	Since int    `json:"since,omitempty"`
	Until   int    `json:"until,omitempty"`
}

// New creates a SeedFaker instance.
func New(opts Options) (*SeedFaker, error) {
	optsJSON, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("seedfaker: marshal opts: %w", err)
	}
	cOpts := C.CString(string(optsJSON))
	defer C.free(unsafe.Pointer(cOpts))

	handle := C.sf_create(cOpts)
	if handle == nil {
		return nil, fmt.Errorf("seedfaker: sf_create: %s", C.GoString(C.sf_last_error()))
	}
	f := &SeedFaker{handle: handle}
	runtime.SetFinalizer(f, func(f *SeedFaker) { f.Close() })
	return f, nil
}

// Close releases the native handle.
func (f *SeedFaker) Close() {
	if f.handle != nil {
		C.sf_destroy(f.handle)
		f.handle = nil
	}
}

// Generate returns a single field value.
func (f *SeedFaker) Generate(field string) (string, error) {
	cField := C.CString(field)
	defer C.free(unsafe.Pointer(cField))

	ptr := C.sf_generate(f.handle, cField)
	if ptr == nil {
		return "", fmt.Errorf("seedfaker: sf_generate(%s): %s", field, C.GoString(C.sf_last_error()))
	}
	val := C.GoString(ptr)
	C.sf_free(ptr)
	return val, nil
}

// RecordOptions for generating multiple records.
type RecordOptions struct {
	Fields  []string `json:"fields"`
	N       int      `json:"n"`
	Ctx     string   `json:"ctx,omitempty"`
	Corrupt string   `json:"corrupt,omitempty"`
}

// GenerateRecords returns records as []map[string]string.
func (f *SeedFaker) GenerateRecords(opts RecordOptions) ([]map[string]string, error) {
	optsJSON, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("seedfaker: marshal opts: %w", err)
	}
	cOpts := C.CString(string(optsJSON))
	defer C.free(unsafe.Pointer(cOpts))

	ptr := C.sf_generate_bulk(f.handle, cOpts)
	if ptr == nil {
		return nil, fmt.Errorf("seedfaker: sf_generate_bulk: %s", C.GoString(C.sf_last_error()))
	}
	jsonStr := C.GoString(ptr)
	C.sf_free(ptr)

	var records []map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &records); err != nil {
		return nil, fmt.Errorf("seedfaker: unmarshal records: %w", err)
	}
	return records, nil
}

// Fields returns all available field names.
func Fields() ([]string, error) {
	ptr := C.sf_fields_json()
	if ptr == nil {
		return nil, fmt.Errorf("seedfaker: sf_fields_json failed")
	}
	jsonStr := C.GoString(ptr)
	C.sf_free(ptr)

	var fields []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &fields); err != nil {
		return nil, err
	}
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names, nil
}

// Fingerprint returns the algorithm version identifier.
func Fingerprint() (string, error) {
	ptr := C.sf_fingerprint()
	if ptr == nil {
		return "", fmt.Errorf("seedfaker: sf_fingerprint failed")
	}
	val := C.GoString(ptr)
	C.sf_free(ptr)
	return val, nil
}
