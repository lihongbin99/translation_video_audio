chrome.runtime.onMessage.addListener((msg, sender) => {
  // 处理 timedtext 响应
  if (msg?.type === "TIMEDTEXT_RESPONSE") {
    const data = msg.payload;
    // chrome.storage.local.set({ timedText: data });
    try {
      const body = JSON.parse(data.body);
      fetch("http://localhost:13520/text", {
        method: "POST",
        body: JSON.stringify(body),
      });
    } catch { }
  }
  
  // 处理当前视频时间
  if (msg?.type === "CURRENT_VIDEO_TIME") {
    const time = msg.payload.time;
    fetch(`http://localhost:13520/time?time=${time}`, {
      method: "GET",
    });
  }
});
