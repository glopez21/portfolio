// https://docs.astro.build/en/guides/content-collections/
import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

export const pillars = ['cybersecurity', 'ai-ml', 'python', 'rust'] as const;
export type Pillar = (typeof pillars)[number];

export const pillarMeta: Record<Pillar, { label: string; blurb: string }> = {
	cybersecurity: {
		label: 'Cybersecurity',
		blurb: 'Offensive & defensive engineering — SIEM/SOAR, evasion research, DFIR, firewall appliances.',
	},
	'ai-ml': {
		label: 'AI / ML',
		blurb: 'From-scratch agents, RAG, evals, and inference infrastructure — no black boxes.',
	},
	python: {
		label: 'Python',
		blurb: 'Well-engineered applications and data tooling built on modern, typed Python.',
	},
	rust: {
		label: 'Rust',
		blurb: 'Systems work from scratch — zero-copy capture, eBPF, engines, and CLIs.',
	},
};

const projects = defineCollection({
	loader: glob({ pattern: '**/*.md', base: './src/content/projects' }),
	schema: z.object({
		title: z.string(),
		pillar: z.enum(pillars),
		tagline: z.string(),
		status: z.enum(['live', 'dev', 'planned', 'research']),
		order: z.number().optional(),
		featured: z.boolean().default(false),
		path: z.string().optional(),
		repo: z.string().url().optional(),
		live: z.string().url().optional(),
		stack: z.array(z.string()).default([]),
		highlights: z.array(z.string()).default([]),
		draft: z.boolean().default(false),
	}),
});

export const collections = { projects };