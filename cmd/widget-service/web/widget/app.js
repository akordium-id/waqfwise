(function () {
  "use strict";

  var CFG = window.WIDGET_CONFIG || {};
  var API_BASE = CFG.apiBase || "/api/v1/widget";

  function formatNumber(n) {
    return new Intl.NumberFormat("id-ID").format(n);
  }

  function formatType(type) {
    var map = {
      land: "Tanah",
      building: "Bangunan",
      cash: "Kas",
      education: "Pendidikan",
      healthcare: "Kesehatan",
      general: "Umum",
    };
    return map[type] || type;
  }

  function escapeHtml(str) {
    var div = document.createElement("div");
    div.appendChild(document.createTextNode(str || ""));
    return div.innerHTML;
  }

  function show(id) {
    document.getElementById(id).classList.remove("hidden");
  }
  function hide(id) {
    document.getElementById(id).classList.add("hidden");
  }

  function renderCampaigns(campaigns) {
    if (!campaigns.length) {
      hide("skeleton");
      show("empty");
      return;
    }

    var grid = document.getElementById("grid");
    grid.innerHTML = campaigns
      .map(function (c) {
        return (
          '<article class="waqf-card" data-url="' +
          escapeHtml(c.donate_url) +
          '">' +
          (c.image_url
            ? '<img class="waqf-card-img" src="' +
              escapeHtml(c.image_url) +
              '" alt="" loading="lazy">'
            : '<div class="waqf-card-img-placeholder"></div>') +
          '<div class="waqf-card-body">' +
          '<span class="waqf-badge waqf-badge-type">' +
          formatType(c.type) +
          "</span>" +
          (c.is_urgent
            ? '<span class="waqf-badge waqf-badge-urgent">Urgent</span>'
            : "") +
          '<h3 class="waqf-card-title">' +
          escapeHtml(c.title) +
          "</h3>" +
          '<p class="waqf-card-desc">' +
          escapeHtml(c.short_desc || "") +
          "</p>" +
          '<div class="waqf-progress"><div class="waqf-progress-bar" style="width:' +
          Math.min(c.progress || 0, 100) +
          '%"></div></div>' +
          '<div class="waqf-meta"><span>Rp ' +
          formatNumber(c.current_amount) +
          " / Rp " +
          formatNumber(c.goal_amount) +
          "</span><span>" +
          formatNumber(c.donor_count) +
          " donatur</span></div>" +
          '<button class="waqf-btn">Donasi Sekarang</button>' +
          "</div></article>"
        );
      })
      .join("");

    grid.addEventListener("click", function (e) {
      var card = e.target.closest(".waqf-card");
      if (card && card.dataset.url) {
        window.top.location.href = card.dataset.url;
      }
    });

    hide("skeleton");
    show("grid");
  }

  function loadCampaigns() {
    hide("error");
    hide("empty");
    hide("grid");
    show("skeleton");

    var params = new URLSearchParams({
      mosque_slug: CFG.mosqueId || "",
      limit: CFG.limit || "6",
      status: "active",
    });
    if (CFG.types) {
      params.set("types", CFG.types);
    }

    fetch(API_BASE + "/campaigns?" + params.toString())
      .then(function (res) {
        if (!res.ok) throw new Error("API error: " + res.status);
        return res.json();
      })
      .then(function (data) {
        renderCampaigns(data.campaigns || []);
      })
      .catch(function () {
        hide("skeleton");
        show("error");
      });
  }

  window.loadCampaigns = loadCampaigns;
  loadCampaigns();
})();
