import pandas as pd
import json
import seaborn as sns
import matplotlib.pyplot as plt

# بارگذاری داده‌ها
df = pd.read_csv("metrics.csv")

# نرمال‌سازی project_id به project1, project2, ...
proj_map = {
    pid: f"project{idx+1}"
    for idx, pid in enumerate(df["project_id"].unique())
}
df["project"] = df["project_id"].map(proj_map)

# خواندن و دسته‌بندی searchable_keys
def categorize_keys(s):
    d = json.loads(s)
    if not d:
        return "none"
    return ",".join(f"{k}={v}" for k, v in d.items())

df["keys_cat"] = df["searchable_keys"].apply(categorize_keys)

# محاسبه نسبت آیتم‌های بازگشتی به کل
df["return_ratio"] = df["returned_items"] / df["total_count"].replace(0, 1)

sns.set(style="whitegrid", palette="muted", font_scale=1.1)


# نمودار توزیع زمان پاسخ (Boxplot)پ
plt.figure(figsize=(8, 5))
sns.boxplot(
    data=df,
    x="scenario",
    y="duration_ms",
    hue="project",
    showfliers=True
)
plt.title("Response Time Distribution by Scenario and Project")
plt.xlabel("Scenario")
plt.ylabel("Duration (ms)")
plt.legend(title="Project")
plt.tight_layout()
plt.savefig("boxplot_duration.png")
plt.close()


#روند زمان پاسخ بر حسب شماره صفحه (Line Plot)
plt.figure(figsize=(8, 5))
sns.lineplot(
    data=df,
    x="page",
    y="duration_ms",
    hue="scenario",
    style="project",
    markers=True,
    dashes=False
)
plt.title("Response Time vs Page Number")
plt.xlabel("Page")
plt.ylabel("Duration (ms)")
plt.legend(title="Scenario / Project", bbox_to_anchor=(1.05, 1), loc="upper left")
plt.tight_layout()
plt.savefig("line_duration_by_page.png")
plt.close()


#نسبت آیتم‌های بازگشتی (Bar Plot)
plt.figure(figsize=(8, 5))
ax = sns.barplot(
    data=df,
    x="scenario",
    y="return_ratio",
    hue="scenario",
    palette="pastel",
    errorbar="sd"
)

# حذف legend تنها در صورتی که وجود داشته باشد
leg = ax.get_legend()
if leg is not None:
    leg.remove()

plt.title("Average Return Ratio per Scenario")
plt.xlabel("Scenario")
plt.ylabel("Returned / Total Count")
plt.ylim(0, 1)
plt.tight_layout()
plt.savefig("bar_return_ratio.png")
plt.close()




# نقشه حرارتی صفحات پویا (Heatmap of Ratios)
pivot_ratio = df.pivot_table(
    index="page",
    columns="scenario",
    values="return_ratio",
    aggfunc="mean"
)

plt.figure(figsize=(8, 4))
sns.heatmap(
    pivot_ratio,
    annot=True,
    fmt=".2f",
    cmap="RdYlBu_r"
)
plt.title("Heatmap of Return Ratio by Page & Scenario")
plt.xlabel("Scenario")
plt.ylabel("Page")
plt.tight_layout()
plt.savefig("heatmap_return_ratio.png")
plt.close()

#ماتریس جفت‌نمودارها (Pairplot)
cols = ["duration_ms", "returned_items", "total_count", "return_ratio"]
sns.pairplot(
    df[cols + ["scenario"]],
    hue="scenario",
    diag_kind="kde",
    corner=True
)
plt.suptitle("Pairwise Relationships Between Metrics", y=1.02)
plt.savefig("pairplot_metrics.png")
plt.close()


