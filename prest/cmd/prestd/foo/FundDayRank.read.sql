SELECT a.*, b.c17
from rank{{.table}} a
         left join fund{{.table}} b on a.c1 = b.c1 and a.c3 = '{{.date}}' and a.c4 != '--'
order by a.c4::float desc
