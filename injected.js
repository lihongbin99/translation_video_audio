(function () {
  // --- Hook XMLHttpRequest ---
  const origOpen = XMLHttpRequest.prototype.open;
  const origSend = XMLHttpRequest.prototype.send;

  XMLHttpRequest.prototype.open = function (method, url, ...rest) {
    this.__ext_url = url;
    return origOpen.call(this, method, url, ...rest);
  };

  XMLHttpRequest.prototype.send = function (...args) {
    this.addEventListener("load", function () {
      try {
        const url = this.__ext_url || "";
        if (url.includes("/api/timedtext")) {
          // 发给 content.js
          window.postMessage(
            { source: "YouTube Audio Extension", type: "TIMEDTEXT_RESPONSE", url, body: this.responseText, ts: Date.now() },
            "*"
          );
        }
      } catch (e) { }
    });
    return origSend.apply(this, args);
  };

})();