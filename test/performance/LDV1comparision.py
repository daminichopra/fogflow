import seaborn as sns
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt


resp = [[  16 ,20, 31, 59, 132,271], [ 10, 11, 15, 19 ,33, 72 ]]

data = pd.DataFrame(resp).T
data.columns = ['NGSILD', 'NGSIV1']
data = data.melt()
data['Thread'] = [1,10,50,100,200,400, 1,10,50,100,200,400]
data


fig, ax = plt.subplots(1,1, figsize=(10,6))
sns.barplot(data=data, x='variable', y='value', hue='Thread', capsize=0.02)
ax.set_ylabel('Response Time(ms)', fontsize='x-large')
ax.set_xlabel('')
ax.tick_params(labelsize='large')
