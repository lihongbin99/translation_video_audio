// 注入 injected.js
const s = document.createElement("script");
s.src = chrome.runtime.getURL("injected.js");
s.onload = () => s.remove();
(document.documentElement || document.head).appendChild(s);

// 接收 injected.js 发来的数据
window.addEventListener("message", (event) => {
  if (event.source !== window) return;
  const msg = event.data;
  if (!msg || msg.source !== "YouTube Audio Extension" || msg.type !== "TIMEDTEXT_RESPONSE") return;

  // 发给 background.js
  chrome.runtime.sendMessage({
    type: "TIMEDTEXT_RESPONSE",
    payload: msg
  });
});

// 获取当前视频时间
function currentVideoTime() {
  const video = document.querySelector("video");
  if (!video) return 0;

  if (!video.paused && !video.ended && video.currentTime > 0) {
    return video.currentTime;
  }
  return 0;
}

// 循环获取当前视频时间
setInterval(() => {
  const currentTime = currentVideoTime();
  if (currentTime && currentTime > 0) {
    chrome.runtime.sendMessage({
      type: "CURRENT_VIDEO_TIME",
      payload: { time: currentTime }
    });
  }
}, 100);
