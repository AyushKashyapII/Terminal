package main

// SiteDomain is shown in the curl footer hint (set to your production host).
const SiteDomain = "ayushkashyap.me"

// URLs — adjust project/demo links to match your public repos.
const (
	GmailURL = "mailto:ayushkashyap0211@gmail.com"
	LinkedInURL = "https://www.linkedin.com/in/ayush-kashyap-9492422a8/"
	GitHubURL   = "https://github.com/AyushKashyapII"
	TwitterURL  = "https://x.com/AyushKashyapII"
	// Set to your live WASM demo, GitHub repo, or site when ready.
	GoChessDemoURL = "https://chess-dzbr1y3df-ayush-kashyaps-projects-9064ed5b.vercel.app"
	PyGitURL       = "https://github.com/AyushKashyapII/Git_Local"
	PacmanRLURL    = "https://github.com/AyushKashyapII/Pacman"
)

const (
	Email        = "ayushkashyap0211@gmail.com"
	Location     = "Chandigarh, India"
)

const welcomeTitle = "Ayush Kashyap"

const welcomeBody = `Student & software engineer — CLIs, agents, and systems that feel good in a real terminal.
I try to build low latency and agentic systems.`

const aboutText = `I'm Ayush Kashyap, I build backends, agentic workflows,
and performance-minded tools — Go, TypeScript, and whatever fits the problem.`

const introBlurb = `B.Tech Mechanical Engineering @ Punjab Engineering College, Chandigarh (2023–2027, CGPA 8.28).
Focused on LLM orchestration, RAG, game engines, and low-level systems — from WASM chess to Git internals.`

const educationBlock = `  Punjab Engineering College, Chandigarh
  B.Tech — Mechanical Engineering
  2023 – 2027  ·  CGPA 8.28`

const skillsBlock = `  Languages   JavaScript, TypeScript, Go, Python, C++, C#, Java
  Frameworks  React, Next.js, Node.js, Express.js, React Native
  AI / ML     TensorFlow, OpenCV, LangChain, LangGraph, Reinforcement Learning
  Tools       Git, Docker, Firebase, Redis, Postman, Vercel, WebAssembly`

const experienceBlock = `  Software Engineer — Constructure AI (startup, Chicago)
  Dec 2025 – April 2026

  • LLM-driven orchestration engine: decompose complex queries into execution plans (DAG)
    and coordinate agents (BIM, Document, Web).
  • Agentic BIM system: dynamic PostgreSQL over 10k+ building elements; natural language
    queries on spatial / architectural properties.
  • Hybrid RAG: PgVector semantic search + keyword matching + section-aware retrieval
    (~25% retrieval accuracy gain).
  • Data pipelines: IFC/SVF2 ingestion, relational schema mapping, Redis Streams for
    real-time updates (~35% lower response latency).`

const projectsBlock = `  GoChess Engine
    Full chess engine: minimax + alpha-beta (~60% search reduction), Zobrist transposition
    tables, quiescence search, material + piece-square evaluation. WASM build (~3× vs JS),
    multi-worker parallel search (~40% faster move compute).

  PyGit — Git clone
    Blob, tree, commit objects; SHA-1 content-addressable storage; staging index; commit DAG
    traversal; CLI: init, add, commit, log.

  Pacman AI (RL)
    Q-learning / DQN agent; state design, reward shaping (~30% faster convergence),
    ~2× average score; epsilon-greedy exploration tuning.`

