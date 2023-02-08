package tracing

import "fmt"

// KeptnCarrier carries the TraceContext
type KeptnCarrier map[string]interface{}

// Get returns the value associated with the passed key.
func (kc KeptnCarrier) Get(key string) string {
	return fmt.Sprintf("%v", kc[key])
}

// Set stores the key-value pair.
func (kc KeptnCarrier) Set(key string, value string) {
	kc[key] = value
}

// Keys lists the keys stored in this carrier.
func (kc KeptnCarrier) Keys() []string {
	keys := make([]string, 0, len(kc))
	for k := range kc {
		keys = append(keys, k)
	}
	return keys
}
