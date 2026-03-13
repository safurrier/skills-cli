---
name: pdf-helper
description: Extract text and tables from PDF files, fill forms, merge documents, split pages. Use when working with PDF documents or when the user mentions PDFs, forms, or document extraction.
license: Apache-2.0
---

# PDF Processing Guide

See [references/REFERENCE.md](references/REFERENCE.md) for detailed library documentation.

## Quick Start

### Extract text
```python
import pdfplumber
with pdfplumber.open("file.pdf") as pdf:
    for page in pdf.pages:
        print(page.extract_text())
```

### Merge PDFs
```python
from pypdf import PdfWriter
writer = PdfWriter()
for path in ["a.pdf", "b.pdf"]:
    writer.append(path)
writer.write("merged.pdf")
```

## Guidelines

- Always check if a PDF is scanned (image-based) before attempting text extraction
- For scanned PDFs, use OCR: `pytesseract` + `pdf2image`
- For form filling, prefer `pypdf` for AcroForm PDFs
