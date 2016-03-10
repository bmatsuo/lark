#Module fun

##Description

The fun module provides a simple API for basic functional
programming.  The module uses a naming convension to distinguish
multiple flavors of the same function.  The map() and sel()
functions operate on table keys and values while the vmap() and
vsel() functions keep the same table keys and only transform table
values.

Most functions operate either on array data or tabular name-value
data unless their documentation explicitly states otherwise (e.g.
the flatten() function).

##Functions

**[flatten](#function-funflatten)**

Returns flat array containing non-table elements of nested array
values.

**[map](#function-funmap)**

Returns a copy of a table with its key-value pairs transformed by a
given function.

**[sel](#function-funsel)**

**[vmap](#function-funvmap)**

Returns a copy of a table with its values transformed by a given
function.

**[vsel](#function-funvsel)**

Returns an table containing elements with values matched by a given
function.

##Function fun.flatten

###Signature

t => tmap

###Description

Returns flat array containing non-table elements of nested array values.

###Parameters

**t** _array_

An array with possibly nested tables.

**tmap** _array_

A copy of t with nested arrays flattened.

**d** _(optional) number_

A depth at which to stop flattening.  A value of zero will
return a copy of the array t.  A negative value or a nil value
will flatten nested arrays at all depths.

##Function fun.map

###Signature

(t, (k, v) => (kmap, vmap)) => tmap

###Description

Returns a copy of a table with its key-value pairs transformed by a given function.

###Parameters

**t** _table_

A table.

**tmap** _table_

A copy of t with its key-value pairs transformed by the given function.

**k** _key_

A key contained in t.

**v** _key_

The value corresponding to k in t.

**kmap** _(optional) key_

A key to store in tmap corresponding to the input key-value
pair.  If no value is returned then k will be used as a key in
tmap.

**vmap** _(optional) any_

The value to store in tmap at key kmap. If nil, or no value is
returned then kmap is removed from tmap.

##Function fun.sel

##Function fun.vmap

###Signature

(t, (v) => (vmap)) => tmap

###Description

Returns a copy of a table with its values transformed by a given function.

###Parameters

**t** _table_

A table.

**tmap** _table_

A copy of t with its values transformed by the given function.

**v** _any_

A value in t.

**vmap** _(optional) any_

The value to store in tmap corresponding to v.  If nil, or no
value is returned then key which v is associated with in t is
not included in tmap.

##Function fun.vsel

###Signature

(t, (v) => keep) => tsel

###Description

Returns an table containing elements with values matched by a given function.

###Parameters

**t** _table_

A table.

**v** _any_

A value in t.

**keep** _boolean_

If true then the input value will be included in tsel under the
same key it had in t.

**tsel** _table_

A table containing values from t for which keep was true.

