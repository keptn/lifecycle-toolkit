---
title: Analysis
description: Reference information about the Metrics API
---

<!-- TOC -->
- [Analysis](#Analysis)
- [AnalysisList](#AnalysisList)
- [AnalysisSpec](#AnalysisSpec)
- [AnalysisStatus](#AnalysisStatus)
- [ProviderResult](#ProviderResult)
- [Timeframe](#Timeframe)


<a id='Analysis'></a>
## Analysis

Analysis is the Schema for the analyses API

Name     | Doc    | Type                                                                                                        
-------- | --- | ------------------------------------------------------------------------------------------------------------
`metadata` |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta)
`spec    ` |  | [AnalysisSpec](#AnalysisSpec)                                                                               
`status  ` |  | [AnalysisStatus](#AnalysisStatus)                                                                           

<a id='AnalysisList'></a>
## AnalysisList

AnalysisList contains a list of Analysis resources

Name     | Doc    | Type                                                                                                    
-------- | --- | --------------------------------------------------------------------------------------------------------
`metadata` |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#listmeta-v1-meta)
`items   ` |  - *mandatory*  | [[]Analysis](#Analysis)                                                                                 

<a id='AnalysisSpec'></a>
## AnalysisSpec

AnalysisSpec defines the desired state of Analysis

Name               | Doc                                                                                                                                                                                                                                                    | Type                   
------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -----------------------
`timeframe         ` | Timeframe specifies the range for the corresponding query in the AnalysisValueTemplate. Please note that either a combination of 'from' and 'to' or the 'recent' property may be set. If neither is set, the Analysis can not be added to the cluster. - *mandatory*  | [Timeframe](#Timeframe)
`args              ` | Args corresponds to a map of key/value pairs that can be used to substitute placeholders in the AnalysisValueTemplate query. i.e. for args foo:bar the query could be "query:percentile(95)?scope=tag(my_foo_label:{{.foo}})".                         | map[string]string      
`analysisDefinition` | AnalysisDefinition refers to the AnalysisDefinition, a CRD that stores the AnalysisValuesTemplates                                                                                                                                                     - *mandatory*  | ObjectReference        

<a id='AnalysisStatus'></a>
## AnalysisStatus

AnalysisStatus stores the status of the overall analysis returns also pass or warnings

Name         | Doc                                                                                     | Type                                        
------------ | --------------------------------------------------------------------------------------- | --------------------------------------------
`timeframe   ` | Timeframe describes the time frame which is evaluated by the Analysis                   - *mandatory*  | [Timeframe](#Timeframe)                     
`raw         ` | Raw contains the raw result of the SLO computation                                      | string                                      
`pass        ` | Pass returns whether the SLO is satisfied                                               | bool                                        
`warning     ` | Warning returns whether the analysis returned a warning                                 | bool                                        
`state       ` | State describes the current state of the Analysis (Pending/Progressing/Completed)       - *mandatory*  | AnalysisState                               
`storedValues` | StoredValues contains all analysis values that have already been retrieved successfully | [map[string]ProviderResult](#ProviderResult)

<a id='ProviderResult'></a>
## ProviderResult

ProviderResult stores reference of already collected provider query associated to its objective template

Name               | Doc                                                           | Type           
------------------ | ------------------------------------------------------------- | ---------------
`objectiveReference` | Objective store reference to corresponding objective template | ObjectReference
`query             ` | Query represents the executed query                           | string         
`value             ` | Value is the value the provider returned                      | string         
`errMsg            ` | ErrMsg stores any possible error at retrieval time            | string         

<a id='Timeframe'></a>
## Timeframe

Name   | Doc                                                                                                                                   | Type                                                                                            
------ | ------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------
`from  ` | From is the time of start for the query. This field follows RFC3339 time format                                                       | [metav1.Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#time-v1-meta)
`to    ` | To is the time of end for the query. This field follows RFC3339 time format                                                           | [metav1.Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#time-v1-meta)
`recent` | Recent describes a recent timeframe using a duration string. E.g. Setting this to '5m' provides an Analysis for the last five minutes | metav1.Duration                                                                                 

