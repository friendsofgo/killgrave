import{_ as t,r as e,o,c as p,a as s,b as n,d as c,e as i}from"./app-lnPknyyT.js";const u={},l=s("h1",{id:"use-json-schema",tabindex:"-1"},[s("a",{class:"header-anchor",href:"#use-json-schema","aria-hidden":"true"},"#"),n(" Use JSON Schema")],-1),r={href:"https://json-schema.org/",target:"_blank",rel:"noopener noreferrer"},d=i(`<p>To do that you will need to define our <code>json schema</code> first:</p><p><em>e.g. <code>imposters/schemas/create_gopher_request.json</code></em></p><div class="language-json line-numbers-mode" data-ext="json"><pre class="language-json"><code><span class="token punctuation">{</span>
    <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;object&quot;</span><span class="token punctuation">,</span>
    <span class="token property">&quot;properties&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
        <span class="token property">&quot;data&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
            <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;object&quot;</span><span class="token punctuation">,</span>
            <span class="token property">&quot;properties&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
                <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
                    <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;string&quot;</span><span class="token punctuation">,</span>
                    <span class="token property">&quot;enum&quot;</span><span class="token operator">:</span> <span class="token punctuation">[</span>
                        <span class="token string">&quot;gophers&quot;</span>
                    <span class="token punctuation">]</span>
                <span class="token punctuation">}</span><span class="token punctuation">,</span>
                <span class="token property">&quot;attributes&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
                    <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;object&quot;</span><span class="token punctuation">,</span>
                    <span class="token property">&quot;properties&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
                        <span class="token property">&quot;name&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
                            <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;string&quot;</span>
                        <span class="token punctuation">}</span><span class="token punctuation">,</span>
                        <span class="token property">&quot;color&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
                            <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;string&quot;</span>
                        <span class="token punctuation">}</span><span class="token punctuation">,</span>
                        <span class="token property">&quot;age&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
                            <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;integer&quot;</span>
                        <span class="token punctuation">}</span>
                    <span class="token punctuation">}</span><span class="token punctuation">,</span>
                    <span class="token property">&quot;required&quot;</span><span class="token operator">:</span> <span class="token punctuation">[</span>
                        <span class="token string">&quot;name&quot;</span><span class="token punctuation">,</span>
                        <span class="token string">&quot;color&quot;</span><span class="token punctuation">,</span>
                        <span class="token string">&quot;age&quot;</span>
                    <span class="token punctuation">]</span>
                <span class="token punctuation">}</span>
            <span class="token punctuation">}</span><span class="token punctuation">,</span>
            <span class="token property">&quot;required&quot;</span><span class="token operator">:</span> <span class="token punctuation">[</span>
                <span class="token string">&quot;type&quot;</span><span class="token punctuation">,</span>
                <span class="token string">&quot;attributes&quot;</span>
            <span class="token punctuation">]</span>
        <span class="token punctuation">}</span>
    <span class="token punctuation">}</span><span class="token punctuation">,</span>
    <span class="token property">&quot;required&quot;</span><span class="token operator">:</span> <span class="token punctuation">[</span>
        <span class="token string">&quot;data&quot;</span>
    <span class="token punctuation">]</span>
<span class="token punctuation">}</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>With this <code>json schema</code>, we expect a <code>request</code> like this:</p><div class="language-json line-numbers-mode" data-ext="json"><pre class="language-json"><code><span class="token punctuation">{</span>
    <span class="token property">&quot;data&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
        <span class="token property">&quot;type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;gophers&quot;</span><span class="token punctuation">,</span>
        <span class="token property">&quot;attributes&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
            <span class="token property">&quot;name&quot;</span><span class="token operator">:</span> <span class="token string">&quot;Zebediah&quot;</span><span class="token punctuation">,</span>
            <span class="token property">&quot;color&quot;</span><span class="token operator">:</span> <span class="token string">&quot;Purples&quot;</span><span class="token punctuation">,</span>
            <span class="token property">&quot;age&quot;</span><span class="token operator">:</span> <span class="token number">55</span>
        <span class="token punctuation">}</span>
    <span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>Then the imposter could be configured as follows:</p><div class="language-json line-numbers-mode" data-ext="json"><pre class="language-json"><code><span class="token punctuation">[</span>
  <span class="token punctuation">{</span>
    <span class="token property">&quot;request&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
        <span class="token property">&quot;method&quot;</span><span class="token operator">:</span> <span class="token string">&quot;POST&quot;</span><span class="token punctuation">,</span>
        <span class="token property">&quot;endpoint&quot;</span><span class="token operator">:</span> <span class="token string">&quot;/gophers&quot;</span><span class="token punctuation">,</span>
        <span class="token property">&quot;schemaFile&quot;</span><span class="token operator">:</span> <span class="token string">&quot;schemas/create_gopher_request.json&quot;</span><span class="token punctuation">,</span>
        <span class="token property">&quot;headers&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
            <span class="token property">&quot;Content-Type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;application/json&quot;</span>
        <span class="token punctuation">}</span>
    <span class="token punctuation">}</span><span class="token punctuation">,</span>
    <span class="token property">&quot;response&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
        <span class="token property">&quot;status&quot;</span><span class="token operator">:</span> <span class="token number">201</span><span class="token punctuation">,</span>
        <span class="token property">&quot;headers&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
            <span class="token property">&quot;Content-Type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;application/json&quot;</span>
        <span class="token punctuation">}</span>
    <span class="token punctuation">}</span>
  <span class="token punctuation">}</span>
<span class="token punctuation">]</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>The path where the schema is located is relative to where the imposters are.</p>`,8);function k(v,q){const a=e("ExternalLinkIcon");return o(),p("div",null,[l,s("p",null,[n("Sometimes you may need to validate requests more thoroughly. In such case, you can create an imposter that only matches with a valid "),s("a",r,[n("json schema"),c(a)]),n(".")]),d])}const b=t(u,[["render",k],["__file","ht-json.html.vue"]]);export{b as default};
