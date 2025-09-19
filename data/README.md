# Data Directory 📊

This directory holds the Excel files that feed wisdom into Haoma's carnival.

## Expected Files

### SCENARIOS.xlsx
Contains category definitions with columns:
- **Column A**: Category Name (e.g., "Cryptography", "Network Security")
- **Column B**: Description (brief explanation of the category)
- **Column C**: PhDT Flag ("PhDT" or "YES" for binary questions, anything else for regular)

### questions.xlsx
Contains the actual questions with columns:
- **Column A**: Category Name (must match SCENARIOS.xlsx)
- **Column B**: Question Text
- **Column C**: Option A
- **Column D**: Option B  
- **Column E**: Option C (leave empty for PhDT questions)
- **Column F**: Option D (leave empty for PhDT questions)
- **Column G**: Correct Answer ("A", "B", "C", or "D")

## Usage

Place your Excel files here and run:
```bash
make seed-excel
```

If the files are missing, the system will create sample data instead:
```bash
make seed
```

## PhDT Questions

For Phishing Detection (PhDT) questions:
- Mark the category as "PhDT" in SCENARIOS.xlsx
- Leave Option C and Option D empty in questions.xlsx
- Use only "A" or "B" as correct answers

This enforces the binary YES/NO nature of phishing detection questions.
