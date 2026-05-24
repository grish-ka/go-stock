import os

def generate_dump(output_file="code_dump.txt"):
    target_extensions = {'.go', '.md'}
    
    with open(output_file, 'w', encoding='utf-8') as outfile:
        for root, _, files in os.walk('.'):
            # Split the path into parts
            parts = root.split(os.sep)
            # Skip hidden directories like .git, but allow '.' and '..'
            if any(part.startswith('.') for part in parts if part not in ('.', '..')):
                continue
                
            for file in files:
                if file == output_file or file == 'dump.py':
                    continue
                    
                ext = os.path.splitext(file)[1]
                if ext in target_extensions:
                    relative_path = os.path.join(root, file)
                    
                    outfile.write(f"\n{'='*40}\n")
                    outfile.write(f"FILE: {relative_path}\n")
                    outfile.write(f"{'='*40}\n\n")
                    
                    try:
                        with open(relative_path, 'r', encoding='utf-8') as infile:
                            outfile.write(infile.read())
                    except Exception as e:
                        outfile.write(f"[Error reading file: {e}]\n")
                    outfile.write("\n")

    print(f" Successfully created {output_file}")

if __name__ == "__main__":
    generate_dump()