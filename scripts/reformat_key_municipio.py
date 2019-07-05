import os
import pandas as pd
import json
import csv as csv

data_migracion = "../data/csv_data/kpi/kpis_gen_nivel_municipio_generado.csv"

output_path = "../data/csv_data/kpi/kpis_gen_nivel_municipio_generado_limpio.csv"

df_output = pd.read_csv(data_migracion,sep=";",error_bad_lines=False, encoding="utf-8", dtype={'CUSEC': object})
df_output['KEY_MUN']=df_output['KEY_MUN'].apply(lambda x: '{0:0>5}'.format(x))

print(df_output)

df_output.to_csv(output_path, sep=";", index=False, line_terminator="\n" , quoting=csv.QUOTE_NONNUMERIC)
