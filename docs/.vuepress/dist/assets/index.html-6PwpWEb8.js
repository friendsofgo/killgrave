import{_ as i,r as t,o,c as r,a as n,b as e,d as l,w as c,e as p}from"./app-NlqJkyee.js";const d={},u=n("h1",{id:"command-line-interface-cli",tabindex:"-1"},[n("a",{class:"header-anchor",href:"#command-line-interface-cli","aria-hidden":"true"},"#"),e(" Command-line interface (CLI)")],-1),v=["src"],h=n("h2",{id:"killgrave",tabindex:"-1"},[n("a",{class:"header-anchor",href:"#killgrave","aria-hidden":"true"},"#"),e(" Killgrave")],-1),m=n("a",{href:"/config"},"config reference",-1),g=n("a",{href:"/guide"},"guide",-1),b=p(`<p>However, you can tune up some of their settings like the host and port where the mock server is listening to, among others, by providing some configuration settings.</p><p>To provide those settings, you can either use the <a href="#available-flags">available CLI flags</a> or use the <code>-config</code> flag to provide the path to a settings file. In such case, you can either use a JSON or YAML configuration file.</p><h3 id="available-flags" tabindex="-1"><a class="header-anchor" href="#available-flags" aria-hidden="true">#</a> Available flags</h3><p>See below the list of available flags:</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code>$ killgrave <span class="token parameter variable">-h</span>

  <span class="token parameter variable">-config</span> string
        path to the configuration <span class="token function">file</span>
  <span class="token parameter variable">-debugger</span>
        run your server with the debugger
  -debugger-addr string
        debugger address <span class="token punctuation">(</span>default <span class="token string">&quot;localhost:3030&quot;</span><span class="token punctuation">)</span>
  <span class="token parameter variable">-host</span> string
        run your server on a different <span class="token function">host</span> <span class="token punctuation">(</span>default <span class="token string">&quot;localhost&quot;</span><span class="token punctuation">)</span>
  <span class="token parameter variable">-imposters</span> string
        directory where imposters are <span class="token builtin class-name">read</span> from <span class="token punctuation">(</span>default <span class="token string">&quot;imposters&quot;</span><span class="token punctuation">)</span>
  <span class="token parameter variable">-port</span> int
        run your server on a different port <span class="token punctuation">(</span>default <span class="token number">3000</span><span class="token punctuation">)</span>
  -proxy-mode string
        proxy mode <span class="token punctuation">(</span>choose between <span class="token string">&#39;all&#39;</span>, <span class="token string">&#39;missing&#39;</span> or <span class="token string">&#39;none&#39;</span><span class="token punctuation">)</span> <span class="token punctuation">(</span>default <span class="token string">&quot;none&quot;</span><span class="token punctuation">)</span>
  -proxy-url string
        proxy url, use it <span class="token keyword">in</span> combination with proxy-mode
  <span class="token parameter variable">-secure</span>
        run your server using TLS <span class="token punctuation">(</span>https<span class="token punctuation">)</span>
  <span class="token parameter variable">-version</span>
        show the version of the application
  <span class="token parameter variable">-watcher</span>
        <span class="token builtin class-name">enable</span> the <span class="token function">file</span> watcher, <span class="token function">which</span> reloads the server on every <span class="token function">file</span> change
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,5);function f(a,k){const s=t("RouterLink");return o(),r("div",null,[u,n("img",{src:a.$withBase("/img/killgrave.png"),alt:"killgrave",style:{"max-width":"130px"}},null,8,v),h,n("p",null,[e("Killgrave is basically a command-line interface (CLI) that can be used with no explicit configuration, but a set of "),l(s,{to:"/config/#imposters"},{default:c(()=>[e("imposters")]),_:1}),e(". Look at the "),m,e(" or "),g,e(" for further details.")]),b])}const x=i(d,[["render",f],["__file","index.html.vue"]]);export{x as default};
