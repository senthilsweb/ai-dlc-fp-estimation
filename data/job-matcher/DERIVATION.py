#!/usr/bin/env python3
"""Derive a TemplrFP dataset from agent-job-matcher's OpenSpec markdown.

Hierarchy: openspec/changes/<name>/ = Intent, "## Bolt N" = Bolt, "- [x]" = Unit of Work.

FP type/complexity do not exist in the source (no story points, hours or
complexity ratings anywhere), so they are inferred here by keyword heuristic.
That is a FIRST PASS FOR HUMAN REVIEW, not a measurement — see the dataset's
description field. The heuristic is deliberately varied rather than uniform:
assigning every task the same weight would make FP a proxy for task count,
which the model explicitly rejects.
"""
import json, re, os, collections, glob

SRC = '/Users/krs/work/agent-job-matcher/openspec/changes'
OUT = '/Users/krs/work/ai-dlc-fp-estimation/data/job-matcher'

TYPE_RULES = [
    (r'schema|model|dataclass|pydantic|types\.ts|jsonresume|shape of', 'ILF'),
    (r'fetch|scrape|external|integration|marketplace|third.?part|arize|otel|telemetry backend', 'EIF'),
    (r'report|render|letter|export|format|output|card|banner|pdf|screenshot', 'EO'),
    (r'\beval\b|evals|test|verify|verification|rubric|assert|suite|conformance', 'EQ'),
    (r'doc|readme|wiki|comment|changelog|guide', 'EQ'),
    (r'api|endpoint|route|fastapi|openapi|mcp|cli|command|server|bridge', 'EI'),
]
HIGH = r'pipeline|orchestrat|correlat|graph|refactor|migrat|architect|end.to.end|containeriz|renumber|root cause|multi-'
LOW = r'^doc|readme|comment|rename|bump|config|env var|\.gitignore|license|badge'


def classify(text):
    t = text.lower()
    ftype = 'EI'
    for pat, ty in TYPE_RULES:
        if re.search(pat, t):
            ftype = ty
            break
    if re.search(HIGH, t):
        cx = 'H'
    elif re.search(LOW, t):
        cx = 'L'
    else:
        cx = 'A'
    return ftype, cx


def clean(text):
    text = re.sub(r'`([^`]*)`', r'\1', text)
    text = re.sub(r'\*\*([^*]*)\*\*', r'\1', text)
    text = re.sub(r'\s+', ' ', text).strip(' .:—-')
    return text[:150]


def parse_tasks_md(path):
    """-> [(bolt_title, bolt_done, [ (task_text, done) ]) ]"""
    lines = open(path, encoding='utf-8').read().split('\n')
    bolts, cur = [], None
    pending = []  # tasks seen before the first bolt heading

    def new_bolt(title, done):
        return {'title': title, 'done': done, 'tasks': []}

    i = 0
    while i < len(lines):
        line = lines[i]
        m = re.match(r'^##\s+(Bolt\s+\d+.*?)\s*$', line)
        if m:
            raw = m.group(1)
            done = '✅' in raw
            title = clean(re.sub(r'✅.*$', '', raw))
            if cur:
                bolts.append(cur)
            cur = new_bolt(title, done)
            if pending:
                cur_pre = new_bolt('Bolt 0 — Groundwork', True)
                cur_pre['tasks'] = pending
                bolts.insert(0, cur_pre)
                pending = []
            i += 1
            continue
        m = re.match(r'^(\s*)- \[([ xX])\]\s*(.*)$', line)
        if m:
            indent, mark, text = m.groups()
            # absorb indented continuation lines
            j = i + 1
            while j < len(lines) and lines[j].strip() and not re.match(r'^\s*- \[|^#', lines[j]) \
                    and len(lines[j]) - len(lines[j].lstrip()) > len(indent):
                text += ' ' + lines[j].strip()
                j += 1
            entry = (clean(text), mark.lower() == 'x')
            (cur['tasks'] if cur else pending).append(entry)
            i = j
            continue
        i += 1
    if cur:
        bolts.append(cur)
    if pending:  # intent had no bolt headings at all
        bolts.append(new_bolt('Bolt 1 — Implementation', all(d for _, d in pending)))
        bolts[-1]['tasks'] = pending
    return [b for b in bolts if b['tasks']]


def intent_status(folder):
    p = os.path.join(SRC, folder, 'proposal.md')
    if not os.path.exists(p):
        return 'roadmap'
    head = open(p, encoding='utf-8').read()[:600].upper()
    if 'IMPLEMENTED' in head:
        return 'completed'
    if 'APPROVED' in head:
        return 'beta'
    return 'roadmap'


def title_of(folder):
    return folder.replace('-', ' ').title().replace('Cli', 'CLI').replace('Ci', 'CI').replace('Ui', 'UI').replace('Api', 'API')


os.makedirs(OUT, exist_ok=True)
intents = sorted(d for d in os.listdir(SRC) if os.path.isdir(os.path.join(SRC, d)))
products, totals = [], collections.Counter()

for idx, folder in enumerate(intents, 1):
    tpath = os.path.join(SRC, folder, 'tasks.md')
    if not os.path.exists(tpath):
        continue
    bolts = parse_tasks_md(tpath)
    st = intent_status(folder)
    features = []
    for b in bolts:
        caps = []
        for text, done in b['tasks']:
            ftype, cx = classify(text)
            caps.append({
                'name': text or 'Unit of work',
                'type': ftype,
                'complexity': cx,
                'status': 'completed' if done else 'in-progress',
            })
            totals['tasks'] += 1
            totals['done' if done else 'open'] += 1
        features.append({'name': b['title'], 'capabilities': caps})
        totals['bolts'] += 1
    dataFile = f'i{idx}-{folder}.json'
    json.dump({'id': folder, 'name': title_of(folder), 'status': st, 'features': features},
              open(os.path.join(OUT, dataFile), 'w'), indent=2)
    products.append({'shortCode': folder.upper()[:14], 'description': f'AI-DLC Intent: {folder}', 'dataFile': dataFile})
    totals['intents'] += 1

print(f"intents={totals['intents']} bolts={totals['bolts']} tasks={totals['tasks']} "
      f"(done={totals['done']} open={totals['open']})")
json.dump(products, open('/tmp/jm_products.json', 'w'))
