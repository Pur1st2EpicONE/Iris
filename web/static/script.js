const API_BASE = "/api/v1";

function showToast(message, type = "success") {
  const toast = document.createElement("div");
  toast.style.cssText = `position:fixed;bottom:20px;right:20px;padding:14px 24px;border-radius:10px;color:white;font-weight:500;z-index:9999;${type === "success" ? "background:#ec4899" : "background:#ef4444"}`;
  toast.textContent = message;
  document.body.appendChild(toast);
  setTimeout(() => toast.remove(), 3000);
}

function formatDateTime(iso) {
  if (!iso) return "-";
  const d = new Date(iso);
  const day = String(d.getDate()).padStart(2, "0");
  const month = String(d.getMonth() + 1).padStart(2, "0");
  const year = d.getFullYear();
  const hours = String(d.getHours()).padStart(2, "0");
  const minutes = String(d.getMinutes()).padStart(2, "0");
  const seconds = String(d.getSeconds()).padStart(2, "0");
  return `${day}.${month}.${year} ${hours}:${minutes}:${seconds}`;
}

function copyLink(short) {
  navigator.clipboard
    .writeText(`${window.location.origin}${API_BASE}/s/${short}`)
    .then(() => showToast("Copied!"));
}

async function shortenURL(original, alias = "") {
  const payload = { original_url: original };
  if (alias) payload.alias = alias;
  const res = await fetch(`${API_BASE}/shorten`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  const data = await res.json();
  if (!res.ok)
    throw new Error(data.error || data.message || "Failed to shorten");
  if (data.short_code) return data.short_code;
  if (data.short_url) return data.short_url.split("/").pop();
  if (data.result) {
    if (typeof data.result === "string") return data.result;
    if (data.result.short_code) return data.result.short_code;
  }
  throw new Error("Invalid server response");
}

async function getAnalytics(shortCode, groupBy = "") {
  let url = `${API_BASE}/analytics/${encodeURIComponent(shortCode)}`;
  if (groupBy) url += `?group_by=${groupBy}`;
  const res = await fetch(url);
  const json = await res.json();
  if (!res.ok) throw new Error(json.error || "Link not found");
  const result = json.result || json;
  const rawData = result.data || [];
  let visits = [];
  const isGrouped = !!groupBy;
  if (isGrouped) {
    const grouped = {};
    for (const item of rawData) {
      const key = item.key || "unknown";
      if (!grouped[key]) grouped[key] = { key, count: 0 };
      grouped[key].count += item.count || 1;
    }
    visits = Object.values(grouped);
    if (groupBy === "day" || groupBy === "month")
      visits.sort((a, b) => a.key.localeCompare(b.key));
    else if (groupBy === "user_agent") visits.sort((a, b) => b.count - a.count);
  } else {
    visits = rawData.map((v) => ({
      count: v.count || 1,
      time: v.time,
      user_agent: v.user_agent || "Unknown",
    }));
  }
  return { total: result.count || rawData.length, visits, isGrouped, groupBy };
}

async function handleShorten(e) {
  e.preventDefault();
  const original = document.getElementById("original-url").value.trim();
  const alias = document.getElementById("alias").value.trim();
  try {
    const shortCode = await shortenURL(original, alias);
    const full = `${window.location.origin}${API_BASE}/s/${shortCode}`;
    const display = document.getElementById("short-link-display");
    display.innerHTML = `<span style="color:#ec4899">${window.location.origin}${API_BASE}/s/</span>${shortCode}`;
    display.dataset.fullUrl = full;
    document.getElementById("visit-link").href = full;
    document.getElementById("shorten-result").classList.remove("hidden");
    document.getElementById("original-url").value = "";
    document.getElementById("alias").value = "";
    showToast("Link shortened successfully!");
  } catch (err) {
    showToast(err.message, "error");
  }
}

function copyShortLink() {
  const display = document.getElementById("short-link-display");
  if (display.dataset.fullUrl)
    navigator.clipboard
      .writeText(display.dataset.fullUrl)
      .then(() => showToast("Copied!"));
}

function viewAnalyticsFromResult() {
  const display = document.getElementById("short-link-display");
  if (!display.dataset.fullUrl) return;
  const shortCode = display.dataset.fullUrl.split("/s/").pop();
  openAnalyticsModal(shortCode);
}

function renderTable(analytics, tbodyElement, theadRowElement) {
  if (analytics.isGrouped) {
    const label =
      analytics.groupBy === "user_agent"
        ? "User Agent"
        : analytics.groupBy === "day"
          ? "Day"
          : "Month";
    theadRowElement.innerHTML = `<th>Visits</th><th>${label}</th>`;
    tbodyElement.innerHTML = analytics.visits.length
      ? analytics.visits
          .map(
            (v) =>
              `<tr><td><strong>${v.count}</strong></td><td style="font-family:monospace;word-break:break-all;">${v.key}</td></tr>`,
          )
          .join("")
      : `<tr><td colspan="2" style="text-align:center;padding:40px 20px;color:#9ca3af">No data</td></tr>`;
  } else {
    theadRowElement.innerHTML = `<th>Visits</th><th>Time</th><th>User Agent</th>`;
    tbodyElement.innerHTML = analytics.visits.length
      ? analytics.visits
          .map(
            (v) =>
              `<tr><td>${v.count}</td><td>${formatDateTime(v.time)}</td><td style="font-family:monospace;font-size:0.9rem;word-break:break-all;">${v.user_agent}</td></tr>`,
          )
          .join("")
      : `<tr><td colspan="3" style="text-align:center;padding:40px 20px;color:#9ca3af">No data</td></tr>`;
  }
}

async function openAnalyticsModal(shortCode, groupBy = "") {
  const modal = document.getElementById("modal");
  document.getElementById("modal-short").textContent = shortCode;
  try {
    const analytics = await getAnalytics(shortCode, groupBy);
    document.getElementById("modal-total").textContent = analytics.total;
    const tbody = document.getElementById("modal-table-body");
    const theadRow = modal.querySelector("thead tr");
    renderTable(analytics, tbody, theadRow);
    modal.classList.remove("hidden");
  } catch (e) {
    showToast(e.message || "Failed to load analytics", "error");
  }
}

function hideModal() {
  document.getElementById("modal").classList.add("hidden");
}

async function fetchAnalyticsDirect() {
  const input = document.getElementById("analytics-short-input").value.trim();
  const groupBy = document.getElementById("analytics-group-select").value;
  if (!input) return showToast("Enter short code", "error");
  try {
    const analytics = await getAnalytics(input, groupBy);
    document.getElementById("total-visits").textContent = analytics.total;
    const tbody = document.getElementById("visits-table-body");
    const theadRow = document.querySelector("#analytics-result thead tr");
    renderTable(analytics, tbody, theadRow);
    document.getElementById("analytics-result").classList.remove("hidden");
  } catch (e) {
    showToast(e.message || "Failed to load analytics", "error");
  }
}

function closeAnalyticsResult() {
  document.getElementById("analytics-result").classList.add("hidden");
}

function main() {
  document
    .getElementById("shorten-form")
    .addEventListener("submit", handleShorten);
}

document.addEventListener("DOMContentLoaded", main);
