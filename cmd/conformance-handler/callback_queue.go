package main

// callbackQueue accumulates invocation records from callback firings and flushes them on drain.
// Embed in a callback interface struct and use record/indexedRecord to build records.
// Object arguments are held internally until drain registers them in the shared registry.
//
// Typical usage:
//
//	iface.indexedRecord(func(n int) map[string]any {
//	    return map[string]any{
//	        "callback": "btck_SomeCallback",
//	        "entry":    iface.ref(fmt.Sprintf("$%s_%d_btck_SomeCallback_entry", iface.ifaceRef, n), entry),
//	        "height":   42,
//	    }
//	})
type callbackQueue struct {
	records  []map[string]any
	ifaceRef string
	registry *Registry
}

// objectRef marks an object argument to be registered in the shared registry on drain.
// Use callbackQueue.ref to construct one. drain stores obj into the registry as-is, so
// if the kernel passes a view into short-lived memory (e.g. a stack-local), copy it before
// the callback returns and pass the copy. Views into long-lived kernel-owned memory
// (e.g. a btck_BlockTreeEntry valid for the chainstate manager's lifetime) can be passed directly.
type objectRef struct {
	ref      string
	obj      any
	registry *Registry
}

// ref creates an objectRef for an object argument using the deterministic naming pattern:
// $<ifaceRef>_<n>_<callback_typedef>_<arg_name>.
func (q *callbackQueue) ref(name string, obj any) objectRef {
	return objectRef{ref: name, obj: obj, registry: q.registry}
}

// record appends an invocation record with no object arguments.
func (q *callbackQueue) record(r map[string]any) {
	q.records = append(q.records, r)
}

// indexedRecord appends an invocation record that needs n to compute ref names.
// n is the 1-based position of this record in the current queue batch.
func (q *callbackQueue) indexedRecord(fn func(n int) map[string]any) {
	q.records = append(q.records, fn(len(q.records)+1))
}

// drain stores each object argument into the shared registry, making it available to the
// runner for follow-up requests, and returns the resolved records in firing order.
func (q *callbackQueue) drain() []map[string]any {
	result := make([]map[string]any, 0, len(q.records))
	for _, record := range q.records {
		resolved := make(map[string]any, len(record))
		for k, v := range record {
			if o, ok := v.(objectRef); ok {
				o.registry.Store(o.ref, o.obj)
				resolved[k] = map[string]any{"ref": o.ref}
			} else {
				resolved[k] = v
			}
		}
		result = append(result, resolved)
	}
	q.records = nil
	return result
}
