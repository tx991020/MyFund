select a.*
from (select *
      from {{.table}}
      where c11 <> '--' and c13 <> '--'  and c10 <> '--' and c7 <> '--' and c8 <> '--' and c9 <> '--') a
         INNER join (select * from {{.table}} where c1 = '{{.fund}}') b
                    on a.c7::float >= b.c7::float and a.c8::float >= b.c8::float and a.c9::float >= b.c9::float and a.c10::float >= b.c10::float  and a.c11::float*a.c13::float   <= b.c11::float*b.c13::float
order by a.c7::float desc limit 10;