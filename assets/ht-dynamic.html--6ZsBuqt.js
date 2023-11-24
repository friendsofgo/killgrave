import{_ as s,o as n,c as a,e}from"./app-lnPknyyT.js";const t={},o=e(`<h1 id="use-dynamic-responses" tabindex="-1"><a class="header-anchor" href="#use-dynamic-responses" aria-hidden="true">#</a> Use dynamic responses</h1><p>Killgrave allows dynamic responses. Using this feature, Killgrave can return different responses on the same endpoint.</p><p>To do this, all imposters need to be sorted from most restrictive to least. Killgrave tries to match the request with each of the imposters in sequence, stopping at the first imposter that matches the request.</p><p>In the following example, there are defined multiple imposters for the <code>POST /gophers</code> endpoint:</p><div class="language-json line-numbers-mode" data-ext="json"><pre class="language-json"><code><span class="token punctuation">[</span>
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
  <span class="token punctuation">}</span><span class="token punctuation">,</span>
  <span class="token punctuation">{</span>
      <span class="token property">&quot;request&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
          <span class="token property">&quot;method&quot;</span><span class="token operator">:</span> <span class="token string">&quot;POST&quot;</span><span class="token punctuation">,</span>
          <span class="token property">&quot;endpoint&quot;</span><span class="token operator">:</span> <span class="token string">&quot;/gophers&quot;</span>
      <span class="token punctuation">}</span><span class="token punctuation">,</span>
      <span class="token property">&quot;response&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
          <span class="token property">&quot;status&quot;</span><span class="token operator">:</span> <span class="token number">400</span><span class="token punctuation">,</span>
          <span class="token property">&quot;headers&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
              <span class="token property">&quot;Content-Type&quot;</span><span class="token operator">:</span> <span class="token string">&quot;application/json&quot;</span>
          <span class="token punctuation">}</span><span class="token punctuation">,</span>
          <span class="token property">&quot;body&quot;</span><span class="token operator">:</span> <span class="token string">&quot;{\\&quot;errors\\&quot;:\\&quot;bad request\\&quot;}&quot;</span>
      <span class="token punctuation">}</span>
  <span class="token punctuation">}</span>
<span class="token punctuation">]</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>Now,</p><ol><li>Let&#39;s say an incoming request does not match the JSON schema specified in the first imposter&#39;s <code>schemaFile</code>.</li><li>Therefore, Killgrave skips this imposter and tries to match the request against the next configured imposter.</li><li>The next configured imposter is much less restrictive, so the request matches and the associated response is returned.</li></ol>`,7),p=[o];function i(r,c){return n(),a("div",null,p)}const u=s(t,[["render",i],["__file","ht-dynamic.html.vue"]]);export{u as default};
