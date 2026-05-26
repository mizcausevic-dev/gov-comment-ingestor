// SPDX-License-Identifier: AGPL-3.0-or-later

package site

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"github.com/mizcausevic-dev/gov-comment-ingestor/internal/service"
)

const style = `
  :root{
    --bg:#070a0f; --panel:#0b1220; --panel2:#0a1426;
    --line:rgba(120,255,170,.18); --line2:rgba(120,255,170,.10);
    --text:#e9f3ff; --muted:rgba(233,243,255,.72); --muted2:rgba(233,243,255,.55);
    --bert:#37ff8b; --bert2:#19c7ff;
    --warn:#ffcc66; --bad:#ff5c7a; --good:#37ff8b; --plum:#b88cff;
    --shadow: 0 18px 60px rgba(0,0,0,.55);
    --radius: 18px;
    --mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    --sans: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial, "Apple Color Emoji", "Segoe UI Emoji";
  }
  *{box-sizing:border-box} html,body{height:100%}
  body{
    margin:0; font-family:var(--sans); color:var(--text);
    background:
      radial-gradient(1200px 600px at 20% -10%, rgba(55,255,139,.18), transparent 60%),
      radial-gradient(900px 520px at 90% 0%, rgba(25,199,255,.16), transparent 55%),
      radial-gradient(1000px 600px at 50% 110%, rgba(55,255,139,.10), transparent 60%),
      linear-gradient(180deg, #05070c 0%, #070a0f 35%, #05070c 100%);
  }
  .grid-bg{
    position:fixed; inset:0; pointer-events:none; opacity:.12; z-index:-1;
    background-image:
      linear-gradient(to right, rgba(55,255,139,.14) 1px, transparent 1px),
      linear-gradient(to bottom, rgba(55,255,139,.10) 1px, transparent 1px);
    background-size: 46px 46px;
    mask-image: radial-gradient(900px 600px at 40% 10%, #000 60%, transparent 100%);
  }
  .wrap{max-width:1280px; margin:0 auto; padding:24px 22px 80px}
  .topbar{
    display:flex; justify-content:space-between; align-items:flex-start; gap:14px;
    border-bottom:1px solid var(--line2); padding-bottom:14px; margin-bottom:22px;
    font-family:var(--mono); font-size:11px; letter-spacing:.16em; color:var(--muted);
    text-transform:uppercase;
  }
  .topbar .left{color:var(--bert)}
  .topbar .right{text-align:right; color:var(--muted)}
  .topbar .right div{margin-bottom:4px}
  .herorow{display:grid; grid-template-columns: 1.5fr .9fr; gap:18px}
  @media (max-width:1000px){.herorow{grid-template-columns:1fr}}
  .hero{
    background: linear-gradient(180deg, rgba(11,18,32,.95), rgba(8,14,26,.92));
    border:1px solid var(--line); border-radius:22px; padding:28px 28px 24px;
    box-shadow: var(--shadow); position:relative; overflow:hidden;
    border-top:2px solid var(--bert2);
  }
  .hero h1{ font-size:64px; line-height:.95; margin:0 0 18px; letter-spacing:-.5px; font-weight:800; }
  @media (max-width:700px){.hero h1{font-size:42px}}
  .hero p{color:var(--muted); font-size:15px; line-height:1.55; max-width:680px; margin:0 0 18px}
  .chiprow{display:flex; flex-wrap:wrap; gap:8px}
  .meta-chip{
    font-family:var(--mono); font-size:11px; color:var(--muted);
    padding:7px 12px; border-radius:999px; border:1px solid var(--line);
    background:rgba(6,10,18,.4);
  }
  .side{display:flex; flex-direction:column; gap:14px}
  .bluf{
    border:1px solid var(--warn); border-left:4px solid var(--warn);
    background: linear-gradient(180deg, rgba(255,204,102,.06), rgba(11,18,32,.92));
    border-radius:14px; padding:16px 18px;
  }
  .bluf .lbl, .corr .lbl{font-family:var(--mono); font-size:10px; letter-spacing:.18em; text-transform:uppercase}
  .bluf .lbl{color:var(--warn)} .corr .lbl{color:var(--bert)}
  .bluf p, .corr p{color:var(--muted); font-size:13.5px; line-height:1.55; margin:6px 0 0}
  .corr{
    border:1px solid var(--bert); border-left:4px solid var(--bert);
    background: linear-gradient(180deg, rgba(55,255,139,.06), rgba(11,18,32,.92));
    border-radius:14px; padding:16px 18px;
  }
  .toolchips{display:flex; flex-wrap:wrap; gap:8px}
  .toolchip{
    font-family:var(--mono); font-size:11px; padding:6px 12px; border-radius:999px;
    border:1px solid currentColor; background:transparent;
  }
  .tc-claude{color:var(--bert2)} .tc-codex{color:var(--warn)} .tc-gpt{color:var(--bert)} .tc-perplex{color:var(--plum)}
  .section{margin-top:34px}
  .sh{
    display:flex; justify-content:space-between; align-items:baseline; gap:14px;
    padding-bottom:10px; border-bottom:1px solid var(--line2); margin-bottom:14px;
  }
  .sh h2{margin:0; font-size:24px; font-weight:600; letter-spacing:-.2px}
  .sh .note{font-family:var(--mono); font-size:11px; color:var(--muted2); letter-spacing:.16em; text-transform:uppercase}
  .kpis{display:grid; grid-template-columns: repeat(6, 1fr); gap:12px}
  @media (max-width:1100px){.kpis{grid-template-columns: repeat(3, 1fr)}} @media (max-width:640px){.kpis{grid-template-columns: repeat(2, 1fr)}}
  .kpi{
    border:1px solid var(--line); border-radius:14px; padding:14px 14px 12px;
    background: linear-gradient(180deg, rgba(11,18,32,.85), rgba(8,14,26,.65));
  }
  .kpi .v{font-family:var(--mono); font-size:26px; font-weight:600; letter-spacing:-.5px}
  .kpi.amber .v{color:var(--warn)} .kpi.cyan .v{color:var(--bert2)} .kpi.green .v{color:var(--bert)} .kpi.plum .v{color:var(--plum)} .kpi.red .v{color:var(--bad)} .kpi.white .v{color:var(--text)}
  .kpi .lbl{font-family:var(--mono); font-size:10px; letter-spacing:.18em; text-transform:uppercase; color:var(--muted); margin-top:6px}
  .kpi .h{font-size:12px; color:var(--muted); line-height:1.45; margin-top:8px}
  .board{display:grid; grid-template-columns: repeat(2,1fr); gap:14px}
  @media (max-width:1000px){.board{grid-template-columns:1fr}}
  .pcard{
    border:1px solid var(--line); border-radius:16px; padding:18px 20px;
    background: linear-gradient(180deg, rgba(11,18,32,.85), rgba(8,14,26,.65));
  }
  .pcard .ptop{display:flex; justify-content:space-between; align-items:center; margin-bottom:8px}
  .pcard .pnum{font-family:var(--mono); font-size:22px; font-weight:600; color:var(--bert)}
  .pcard .ppri{font-family:var(--mono); font-size:10px; padding:5px 10px; border-radius:999px; border:1px solid var(--line); color:var(--bert); letter-spacing:.14em; background:rgba(55,255,139,.06)}
  .pcard h3{margin:6px 0 8px; font-size:19px; font-weight:600}
  .pcard .pdesc{font-size:13.5px; color:var(--muted); line-height:1.55; margin:0 0 12px}
  .pcard ul.check{list-style:none; padding:0; margin:0}
  .pcard ul.check li{display:grid; grid-template-columns:18px 1fr; gap:10px; padding:6px 0; font-size:13.5px; color:var(--muted); line-height:1.45}
  .pcard ul.check li:before{content:""; width:14px; height:14px; border:1px solid var(--line); border-radius:3px; background:rgba(6,10,18,.4); margin-top:3px}
  .ttbl{
    width:100%; border-collapse:separate; border-spacing:0;
    border:1px solid var(--line); border-radius:14px; overflow:hidden;
  }
  .ttbl th, .ttbl td{padding:13px 14px; text-align:left; font-size:13.5px; vertical-align:top}
  .ttbl thead th{
    font-family:var(--mono); font-size:11px; letter-spacing:.16em; text-transform:uppercase;
    color:var(--muted2); border-bottom:1px solid var(--line); background:rgba(11,18,32,.5);
  }
  .ttbl tbody tr:hover{background:rgba(55,255,139,.03)}
  .ttbl td, .ttbl td *{color:var(--muted)}
  .ttbl b{color:var(--text)}
  .st{font-family:var(--mono); font-size:10px; padding:4px 9px; border-radius:6px; letter-spacing:.1em; text-transform:uppercase; border:1px solid currentColor; display:inline-block}
  .st.good{color:var(--bert)} .st.watch{color:var(--warn)} .st.critical{color:var(--bad)}
  footer{
    margin-top:30px; padding-top:14px; border-top:1px dashed var(--line2);
    display:flex; justify-content:space-between; gap:10px; flex-wrap:wrap;
    font-family:var(--mono); font-size:11px; color:var(--muted2); letter-spacing:.08em;
  }
  a{color:var(--bert2); text-decoration:none}
  a:hover{text-decoration:underline}
`

type Page struct {
	Title      string
	Slug       string
	NavSlug    string
	Eyebrow    string
	HeroTitle  string
	HeroBody   string
	RightCards []SideCard
	Body       string
}

type SideCard struct {
	Label string
	Title string
	Body  string
}

func RenderOverview() string {
	payload := service.SamplePayload()
	kpis := []map[string]string{
		{"class": "green", "value": fmt.Sprintf("%d", payload.Summary.AdaptersHealthy), "label": "Healthy adapters", "hint": "Feeds delivering normalized comments on-time"},
		{"class": "amber", "value": fmt.Sprintf("%d", payload.Summary.AdaptersDegraded), "label": "Degraded adapters", "hint": "Need parser or lineage intervention"},
		{"class": "cyan", "value": fmt.Sprintf("%d", payload.Summary.Comments24h), "label": "Comments / 24h", "hint": "Across federal, state, and counsel intake"},
		{"class": "white", "value": fmt.Sprintf("%.1f%%", payload.Summary.NormalizationRate), "label": "Normalization rate", "hint": "Comments with thread-safe normalized output"},
		{"class": "red", "value": fmt.Sprintf("%d", payload.Summary.ThreadingFlags), "label": "Thread flags", "hint": "Duplicates or broken reply chains"},
		{"class": "plum", "value": fmt.Sprintf("%d", payload.Summary.EvidenceGaps), "label": "Evidence gaps", "hint": "Missing attachments or citation lineage"},
	}

	var kpiHTML strings.Builder
	for _, kpi := range kpis {
		kpiHTML.WriteString(fmt.Sprintf(`<div class="kpi %s"><div class="v">%s</div><div class="lbl">%s</div><div class="h">%s</div></div>`,
			kpi["class"], html.EscapeString(kpi["value"]), html.EscapeString(kpi["label"]), html.EscapeString(kpi["hint"])))
	}

	var board strings.Builder
	for index, item := range payload.IngestLane {
		board.WriteString(fmt.Sprintf(`
      <article class="pcard">
        <div class="ptop"><div class="pnum">G-%02d</div><div class="ppri">%s</div></div>
        <h3>%s</h3>
        <p class="pdesc">%s · %s · %d new comments</p>
        <ul class="check">
          <li>%s</li>
          <li>%s</li>
          <li><strong>Owner:</strong> %s</li>
          <li><strong>Next action:</strong> %s</li>
        </ul>
      </article>`,
			index+1,
			html.EscapeString(strings.ToUpper(item.Status)),
			html.EscapeString(item.Queue),
			html.EscapeString(item.Agency),
			html.EscapeString(item.Window),
			item.NewComments,
			html.EscapeString(item.DriftSignal),
			html.EscapeString(item.AttachmentRisk),
			html.EscapeString(item.Owner),
			html.EscapeString(item.NextAction),
		))
	}

	body := fmt.Sprintf(`
  <section class="section">
    <div class="sh"><h2>Overview</h2><div class="note">queue health + review safety</div></div>
    <div class="kpis">%s</div>
  </section>
  <section class="section">
    <div class="sh"><h2>Ingestion pressure board</h2><div class="note">highest-risk dockets first</div></div>
    <div class="board">%s</div>
  </section>`, kpiHTML.String(), board.String())

	return render(Page{
		Title:     "Gov Comment Ingestor",
		Slug:      "/",
		NavSlug:   "/",
		Eyebrow:   "Gov comment ingestion command surface",
		HeroTitle: "Government comment feeds fail when adapter drift, attachment loss, and thread mismatches stay hidden.",
		HeroBody:  "This operator surface makes the intake layer explicit: which sources are healthy, which dockets are lagging, where parser drift is creeping in, and whether evidence lineage is safe enough for legal, policy, and procurement review.",
		RightCards: []SideCard{
			{Label: "Core offer", Title: "Go-native ingest control plane", Body: "Source adapters, normalization trust, and evidence lineage in one operational surface."},
			{Label: "Buyer fit", Title: "GovTech + RegTech", Body: "Public affairs, legal ops, compliance, procurement, and submission teams who need reliable docket intake without spreadsheet drift."},
			{Label: "Execution style", Title: "Review-safe ingestion", Body: "Track source lag, parser drift, and attachment lineage before they contaminate downstream review packets."},
		},
		Body: body,
	})
}

func RenderIngestLane() string {
	payload := service.SamplePayload()
	var rows strings.Builder
	for _, item := range payload.IngestLane {
		rows.WriteString(fmt.Sprintf(`
      <tr>
        <td><b>%s</b><br>%s</td>
        <td>%s</td>
        <td><span class="st %s">%s</span></td>
        <td>%d</td>
        <td>%s</td>
        <td>%s</td>
      </tr>`,
			html.EscapeString(item.Queue),
			html.EscapeString(item.Agency),
			html.EscapeString(item.Window),
			statusClass(item.Status),
			html.EscapeString(item.Status),
			item.NewComments,
			html.EscapeString(item.Owner),
			html.EscapeString(item.NextAction),
		))
	}
	body := fmt.Sprintf(`
  <section class="section">
    <div class="sh"><h2>Ingest lane</h2><div class="note">queue ownership + drift pressure</div></div>
    <table class="ttbl">
      <thead>
        <tr>
          <th>Docket queue</th><th>Review window</th><th>Status</th><th>New comments</th><th>Owner</th><th>Next action</th>
        </tr>
      </thead>
      <tbody>%s</tbody>
    </table>
  </section>`, rows.String())

	return render(Page{
		Title:     "Gov Comment Ingestor — ingest lane",
		Slug:      "/ingest-lane",
		NavSlug:   "/ingest-lane",
		Eyebrow:   "Review-safe queue ownership",
		HeroTitle: "See which dockets are safe to synthesize and which still carry ingestion risk.",
		HeroBody:  "The ingest lane keeps queue volume, adapter drift, attachment gaps, and owner handoff visible before legal or policy reviewers start drafting against bad source state.",
		RightCards: []SideCard{
			{Label: "Signal", Title: "Thread-safe first", Body: "Do not let comment volume hide broken reply chains or mis-bound attachments."},
			{Label: "Pressure", Title: "Review windows", Body: "Each queue maps intake health to the next review or submission deadline."},
			{Label: "Control", Title: "Owner + next action", Body: "Every degraded queue has a named owner and one clear intervention path."},
		},
		Body: body,
	})
}

func RenderSourceAdapters() string {
	payload := service.SamplePayload()
	var rows strings.Builder
	for _, adapter := range payload.SourceAdapters {
		rows.WriteString(fmt.Sprintf(`
      <tr>
        <td><b>%s</b><br>%s</td>
        <td><span class="st %s">%s</span></td>
        <td>%s</td>
        <td>%s</td>
        <td>%s</td>
        <td>%s</td>
      </tr>`,
			html.EscapeString(adapter.Name),
			html.EscapeString(adapter.Feed),
			statusClass(adapter.Status),
			html.EscapeString(adapter.Status),
			html.EscapeString(adapter.Lag),
			html.EscapeString(adapter.ParserVersion),
			html.EscapeString(adapter.CommentCoverage),
			html.EscapeString(adapter.Escalation),
		))
	}
	body := fmt.Sprintf(`
  <section class="section">
    <div class="sh"><h2>Source adapters</h2><div class="note">parser health + lineage trust</div></div>
    <table class="ttbl">
      <thead>
        <tr>
          <th>Adapter</th><th>Status</th><th>Lag</th><th>Parser version</th><th>Coverage</th><th>Escalation</th>
        </tr>
      </thead>
      <tbody>%s</tbody>
    </table>
  </section>`, rows.String())

	return render(Page{
		Title:     "Gov Comment Ingestor — source adapters",
		Slug:      "/source-adapters",
		NavSlug:   "/source-adapters",
		Eyebrow:   "Adapter reliability lane",
		HeroTitle: "Treat the source adapters as production systems, not background plumbing.",
		HeroBody:  "Government comment workflows break quietly when HTML changes, attachment manifests drift, or inbox forwarding rules lose context. This view keeps parser health, coverage, and escalation posture visible.",
		RightCards: []SideCard{
			{Label: "Core offer", Title: "Adapter drift visibility", Body: "A single table for lag, coverage, parser version, and escalation state."},
			{Label: "Buyer fit", Title: "Operator-owned reliability", Body: "For teams that need ingestion proof before they trust dashboard or briefing outputs."},
			{Label: "Execution style", Title: "Evidence lineage", Body: "Attachment policy and checksum posture stay tied to adapter state."},
		},
		Body: body,
	})
}

func RenderVerification() string {
	payload := service.SamplePayload()
	var cards strings.Builder
	for index, gate := range payload.Verification {
		cards.WriteString(fmt.Sprintf(`
      <article class="pcard">
        <div class="ptop"><div class="pnum">V-%02d</div><div class="ppri">%s</div></div>
        <h3>%s</h3>
        <p class="pdesc">%s</p>
        <ul class="check">
          <li><strong>Risk:</strong> %s</li>
          <li><strong>Next action:</strong> %s</li>
        </ul>
      </article>`,
			index+1,
			html.EscapeString(strings.ToUpper(gate.Status)),
			html.EscapeString(gate.Name),
			html.EscapeString(gate.Evidence),
			html.EscapeString(gate.Risk),
			html.EscapeString(gate.NextAction),
		))
	}
	body := fmt.Sprintf(`
  <section class="section">
    <div class="sh"><h2>Verification</h2><div class="note">release gate for normalized intake</div></div>
    <div class="board">%s</div>
  </section>`, cards.String())

	return render(Page{
		Title:     "Gov Comment Ingestor — verification",
		Slug:      "/verification",
		NavSlug:   "/verification",
		Eyebrow:   "Verification posture",
		HeroTitle: "Clear ingestion safety before you clear analyst or counsel review.",
		HeroBody:  "The verification lane ties adapter health, evidence lineage, duplicate-thread suppression, and export contracts into one buyer-readable release gate.",
		RightCards: []SideCard{
			{Label: "Signal", Title: "No blind ingestion", Body: "Stale or duplicated comment streams become review blockers, not quiet footnotes."},
			{Label: "Proof", Title: "Lineage preserved", Body: "Checksums, manifests, and export contracts stay tied to the intake state."},
			{Label: "Buyer value", Title: "Safer synthesis", Body: "Analysts and procurement teams inherit cleaner evidence instead of parsing noise."},
		},
		Body: body,
	})
}

func RenderDocs() string {
	payload := service.SamplePayload()
	body := fmt.Sprintf(`
  <section class="section">
    <div class="sh"><h2>Docs</h2><div class="note">implementation notes</div></div>
    <div class="board">
      <article class="pcard">
        <div class="ptop"><div class="pnum">A</div><div class="ppri">Adapter model</div></div>
        <h3>Feed adapters</h3>
        <p class="pdesc">Six modeled adapters span federal APIs, state HTML, counsel inbox flows, and attachment mirrors. %d are currently healthy.</p>
        <ul class="check">
          <li>Every adapter carries lag, parser version, coverage, and escalation context.</li>
          <li>Attachment policy is explicit so evidence gaps are buyer-readable.</li>
        </ul>
      </article>
      <article class="pcard">
        <div class="ptop"><div class="pnum">B</div><div class="ppri">Queue model</div></div>
        <h3>Docket intake</h3>
        <p class="pdesc">Queue rows tie agency, review window, status, comment volume, owner, and next action into one review-safe object.</p>
        <ul class="check">
          <li>Thread flags block synthesis before narratives double-count themes.</li>
          <li>Owners are explicit so escalation survives handoffs.</li>
        </ul>
      </article>
    </div>
  </section>`, payload.Summary.AdaptersHealthy)

	return render(Page{
		Title:     "Gov Comment Ingestor — docs",
		Slug:      "/docs",
		NavSlug:   "/docs",
		Eyebrow:   "Operator docs",
		HeroTitle: "The ingestion tier is part of the product surface, not just a backend utility.",
		HeroBody:  "This Go repo demonstrates how Kinetic Gain can expose regulated workflow intake as a review-safe operator dashboard with static deployment and deterministic sample data.",
		RightCards: []SideCard{
			{Label: "Language atlas", Title: "Go surface", Body: "This opens the Language Atlas with a Go-native ingestion and prerender path."},
			{Label: "Deploy", Title: "Static Pages + CNAME", Body: "Prerendered HTML and JSON ship to GitHub Pages under the wildcard domain."},
			{Label: "Embedded tie-back", Title: "Review-safe analytics", Body: "The same primitive can sit inside regulated SaaS products without unsafe writes."},
		},
		Body: body,
	})
}

func JSON(v any) string {
	bytes, _ := json.MarshalIndent(v, "", "  ")
	return string(bytes)
}

func render(page Page) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s</title>
  <style>%s</style>
</head>
<body>
  <div class="grid-bg"></div>
  <div class="wrap">
    <div class="topbar">
      <div class="left">KINETIC GAIN · gov comment ingestion</div>
      <div class="right">
        <div>Go operator surface · RegTech / GovTech</div>
        <div>Evidence lineage · parser drift · review-safe intake</div>
      </div>
    </div>
    <div class="herorow">
      <section class="hero">
        <div class="chiprow">
          <span class="meta-chip">%s</span>
          <span class="meta-chip">CNAME · ingestor.kineticgain.com</span>
          <span class="meta-chip">Go · prerendered Pages deploy</span>
        </div>
        <h1>%s</h1>
        <p>%s</p>
        <div class="chiprow">
          %s
        </div>
      </section>
      <aside class="side">%s</aside>
    </div>
    %s
    <footer>
      <div>gov-comment-ingestor · AGPL-3.0-or-later · synthetic demonstration data only</div>
      <div>Routes: / · /ingest-lane · /source-adapters · /verification · /docs</div>
    </footer>
  </div>
</body>
</html>`,
		html.EscapeString(page.Title),
		style,
		html.EscapeString(page.Eyebrow),
		page.HeroTitle,
		html.EscapeString(page.HeroBody),
		nav(page.NavSlug),
		renderSideCards(page.RightCards),
		page.Body,
	)
}

func nav(active string) string {
	items := []struct {
		Label string
		Href  string
	}{
		{"Overview", "/"},
		{"Ingest Lane", "/ingest-lane"},
		{"Source Adapters", "/source-adapters"},
		{"Verification", "/verification"},
		{"Docs", "/docs"},
	}
	var out strings.Builder
	for _, item := range items {
		class := "toolchip tc-claude"
		if item.Href == active {
			class = "toolchip tc-gpt"
		}
		out.WriteString(fmt.Sprintf(`<a class="%s" href="%s">%s</a>`, class, item.Href, html.EscapeString(item.Label)))
	}
	return out.String()
}

func renderSideCards(cards []SideCard) string {
	var out strings.Builder
	for index, card := range cards {
		class := "corr"
		if index == 0 {
			class = "bluf"
		}
		out.WriteString(fmt.Sprintf(`<article class="%s"><div class="lbl">%s</div><p><strong>%s</strong><br>%s</p></article>`,
			class,
			html.EscapeString(card.Label),
			html.EscapeString(card.Title),
			html.EscapeString(card.Body),
		))
	}
	return out.String()
}

func statusClass(status string) string {
	switch strings.ToLower(status) {
	case "healthy":
		return "good"
	case "critical":
		return "critical"
	default:
		return "watch"
	}
}
