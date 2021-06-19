select a.*
from (select *
      from {{.table}}
      where c11 <> '--' and c13 <> '--' and c10 <> '--' and c6 <> '--' and c8 <> '--' and c9 <> '--') a
         INNER join (select * from {{.table}} where c1 = '163402') b
                    on a.c6::float >= b.c6::float and a.c8::float >= 5 and a.c9::float >= 8 and a.c10::float >= 10  and a.c11::float*a.c13::float   <= b.c11::float*b.c13::float
order by a.c6:: float desc