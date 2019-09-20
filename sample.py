import sys, csv, random
from datetime import date, timedelta

persons = [
    ('person1', 'A'),
    ('person2', 'A'),
    ('person3', 'A'),
    ('person4', 'B'),
    ('person5', 'B'),
    ('person6', 'B'),
    ('person7', 'C'),
    ('person8', 'C'),
    ('person9', 'C')
]

def date_range():
    base = date(2018, 11, 1)
    return [base + timedelta(days=x) for x in range(0, 30)]

def write_csv(dt, rows):
    file_path = f'./csv/{date.strftime(dt, "%Y%m%d")}.csv'
    with open(file_path, 'w') as f:
        writer = csv.writer(f)
        writer.writerows(rows)

def make_rows(dt_str):
    p = random.choice(persons)
    rows = [dt_str, p[0], p[1], random.randint(1, 100)]
    return rows

if __name__ == '__main__':
    for dt in date_range():
        dt_str = date.strftime(dt, "%Y-%m-%d")
        rows = [make_rows(dt_str) for i in range(0, 100000000)]
        write_csv(dt, rows)
