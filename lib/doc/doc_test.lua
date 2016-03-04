function test_doc()
	local doc = require('doc')
	local x = function(p) return p .. '.xyz' end

	local sig_string = 'p => string'
	local sig = doc.sig(sig_string)
	sig(x)
	docs = doc.get(x)
	assert(docs.sig == sig_string)
	assert(not docs.desc)
	assert(not docs.params)

	local desc_string = 'return things'
	local desc = doc.desc(desc_string)
	desc(x)
	docs = doc.get(x)
	assert(docs.sig == sig_string)
	assert(docs.desc == desc_string)
	assert(not docs.params)

	local param_string = 'p  a thing'
	local param = doc.param(param_string)
	param(x)
	docs = doc.get(x)
	assert(docs.sig == sig_string)
	assert(docs.desc == desc_string)
	assert(table.concat(docs.params, '\n') == param_string)

	sig_string = '(p, b) => string'
	sig = doc.sig(sig_string)
	_ = sig .. x
	docs = doc.get(x)
	assert(docs.sig == sig_string)
	assert(docs.desc == desc_string)
	assert(table.concat(docs.params, '\n') == param_string)

	desc_string = 'return things'
	desc = doc.desc(desc_string)
	_ = desc .. x
	docs = doc.get(x)
	assert(docs.sig == sig_string)
	assert(docs.desc == desc_string)
	assert(table.concat(docs.params, '\n') == param_string)

	local param_string2 = 'p  a thing'
	param = doc.param(param_string2)
	_ = param .. x
	docs = doc.get(x)
	assert(docs.sig == sig_string)
	assert(docs.desc == desc_string)
	assert(table.concat(docs.params, '\n') == (param_string .. '\n' .. param_string2))
end
