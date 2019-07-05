import os
import pandas as pd
import json
import csv as csv

data_migracion = "../data/csv_data/kpi/kpi_migracion_mun_generado.csv"
data_edad = "../data/csv_data/kpi/kpi_edad_nivel_mun_generado.csv"
data_gen = "../data/csv_data/kpi/kpis_gen_nivel_municipio_generado_limpio.csv"

output_path = "../data/csv_data/kpi/kpis_municipio_combinados.csv"

df_output = pd.read_csv(data_gen,sep=";",error_bad_lines=False, encoding="utf-8", dtype={'KEY_MUN': object})
df_migracion = pd.read_csv(data_migracion,sep=";",error_bad_lines=False, encoding="utf-8", dtype={'KEY_MUN': object})
df_edad = pd.read_csv(data_edad,sep=";",error_bad_lines=False, encoding="utf-8", dtype={'KEY_MUN': object})

df_output["kpi_migracion"] = df_migracion["kpi_migracion"]
df_output["kpi_edad"] = df_edad["kpi_edad"]
df_output["kpi_suma"] = df_output["kpi_migracion"] + df_output["kpi_edad"] + df_output["RANKING_UP"]+ df_output["RANKING_UP_PSOE"]+ df_output["UP_MEDIANA"]

print(df_output)

df_output.to_csv(output_path, sep=";", index=False, line_terminator="\n" , quoting=csv.QUOTE_NONNUMERIC)
