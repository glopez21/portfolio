import type { Pillar } from '../content.config';

export const pillarColor: Record<Pillar, string> = {
	cybersecurity: 'var(--pillar-cybersecurity)',
	'ai-ml': 'var(--pillar-ai-ml)',
	python: 'var(--pillar-python)',
	rust: 'var(--pillar-rust)',
};

export const statusLabel: Record<string, string> = {
	live: 'Live',
	dev: 'In dev',
	planned: 'Planned',
	research: 'Research',
};

export const statusBadge = (status: string) => `badge-status-${status}`;