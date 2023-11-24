import{_ as r,r as t,o as i,c as l,a,b as e,d as s,e as o}from"./app-GKjJbFgT.js";const c={},d=a("h1",{id:"installation",tabindex:"-1"},[a("a",{class:"header-anchor",href:"#installation","aria-hidden":"true"},"#"),e(" Installation")],-1),h={href:"https://semver.org/",target:"_blank",rel:"noopener noreferrer"},p=o(`<p>You can install Killgrave in different ways, but all of them are very simple:</p><h2 id="go-toolchain" tabindex="-1"><a class="header-anchor" href="#go-toolchain" aria-hidden="true">#</a> Go Toolchain</h2><p>One of them is of course using <code>go install</code>, Killgrave is a Go project and therefore can be compiled using the go toolchain:</p><div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code>$ go <span class="token function">install</span> github.com/friendsofgo/killgrave/cmd/killgrave@<span class="token punctuation">{</span>version<span class="token punctuation">}</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><p><em>Note that <code>version</code> must be replaced by the version that you want to install. If left unspecified, the <code>main</code> branch will be installed.</em></p><h2 id="homebrew" tabindex="-1"><a class="header-anchor" href="#homebrew" aria-hidden="true">#</a> Homebrew</h2>`,6),u={href:"https://brew.sh/",target:"_blank",rel:"noopener noreferrer"},b=o(`<div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code>$ brew <span class="token function">install</span> friendsofgo/tap/killgrave
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><p>⚠️ If you are installing via Homebrew, you always get the latest Killgrave version, we hope to fix this soon!</p><h2 id="docker" tabindex="-1"><a class="header-anchor" href="#docker" aria-hidden="true">#</a> Docker</h2>`,3),m={href:"https://www.docker.com/",target:"_blank",rel:"noopener noreferrer"},g=o(`<div class="language-bash line-numbers-mode" data-ext="sh"><pre class="language-bash"><code>$ <span class="token function">docker</span> run <span class="token parameter variable">-it</span> <span class="token parameter variable">--rm</span> <span class="token parameter variable">-p</span> <span class="token number">3000</span>:3000 <span class="token parameter variable">-v</span> <span class="token environment constant">$PWD</span>/:/home <span class="token parameter variable">-w</span> /home friendsofgo/killgrave <span class="token parameter variable">-host</span> <span class="token number">0.0</span>.0.0
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><ul><li><p><code>-p 3000:3000</code> is used to forward the local port <code>3000</code> (Killgrave&#39;s default port) to container&#39;s port <code>3000</code>, otherwise Killgrave won&#39;t be reachable from the host.</p></li><li><p><code>-host 0.0.0.0</code> is used to change the Killgrave&#39;s default host (<code>localhost</code>) to allow Killgrave to listen to and respond to incoming requests from outside the container, otherwise Killgrave won&#39;t be reachable from the host.</p></li></ul><h2 id="other" tabindex="-1"><a class="header-anchor" href="#other" aria-hidden="true">#</a> Other</h2>`,3),v={href:"https://github.com/friendsofgo/killgrave/releases",target:"_blank",rel:"noopener noreferrer"};function f(k,_){const n=t("ExternalLinkIcon");return i(),l("div",null,[d,a("blockquote",null,[a("p",null,[e("⚠️ Even though Killgrave is a very robust tool and is being used by some companies in production environments, it's still in initial development. Therefore, 'minor' version numbers are used to signify breaking changes and 'patch' version numbers are used for non-breaking changes or bugfixing. As soon as v1.0.0 is released, Killgrave will start to use "),a("a",h,[e("SemVer"),s(n)]),e(" as usual.")])]),p,a("p",null,[e("If you are a macOS user, you can install Killgrave using "),a("a",u,[e("Homebrew"),s(n)]),e(":")]),b,a("p",null,[e("Killgrave is also available through "),a("a",m,[e("Docker"),s(n)]),e(".")]),g,a("p",null,[e("Windows and Linux users can download binaries from the "),a("a",v,[e("GitHub Releases"),s(n)]),e(" page.")])])}const x=r(c,[["render",f],["__file","installation.html.vue"]]);export{x as default};
