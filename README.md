# Kvse
## About
A command line interface in go to transform compliant yaml to flat structure of key value pairs.  Kvse-compliant yaml is a subset of yaml requiring a few constraints as follows:

### Kvse Compliant YAML
- ending values can only be scalars
- a list of maps with the key of 'name' has special meaning and will result in data transform

## Transforming Data

There are two general types of transform performed.

### Transform from yaml to key-value pairs without data loss.  

- This is done on any given kvse-compliant yaml input with no special indications needed.

### Schema Key Replacement

 - This is done when the yaml schema includes a list in which the members of the list are maps in which every map includes a key-value pair of `name: <scalar>`.  In this case, kvse is being instructed that the parent element is a metadata key and should be replaced with each of the names provided.

### Why?

Why flatten data from yaml to key/value pairs?

  It has been found in working with hierarchical data, the extraction of information often involves a verbosity of commands and syntax, requiring detailed understanding of the data structure, embedding this code into a somewhat arbitrarily sourced tool of choice, for example, `jq`, `yq`, `jsonpath`, `jq-action`, etc..,  each an added dependency requiring their own upkeep and understanding in addition to the data at hand.  Flattening the data allows one to only need to know a key to get at the value: no logic no coding, no toolchain sprawl.  A flat key-value structure of configuration data, if not overwhelmingly preferred, is at least extremely popular, and describes many configuration stores old and new. 
 
  "But my yaml data is flat already, so I can address it as key/value pairs now.  Why do I need kvse?"
 
 The answer gets at the question of why use yaml in the first place?
  There are probably a lot of reasons, but, most among them one would guess is we are familiar with yaml and it gives us a sense of structure without thinking too much about it. It looks like a `standard` by just typing it, not to mention it  adorns us with instant hero status.

**But working with yaml sources in configuration scripts, automation pipelines and the like - things you don't normally associate with full blown languages, but rather DSLs and the like, having to extract values from json and yaml begins to wear thin.**

Kvse is created to keep our data store in yaml, but consume it as a KVDB. 
With no extensions or anything beyond standard yaml syntax, we can enhance the yaml data, making it more descriptive, more flexible, more understandable, and finally, and perhaps at least as important as all the other '-ables',  more usable.  You can forget downloading and installing `yq, jq, jsonpath, jmespath, objectpath, jq-action, etc...` in your scripts and pipelines and trying to remember all the intricacies of each of these separate query tools in doing so.

"Hey! It sounds like you are suggesting there *is* something wrong with my yaml after all - what gives?".  OK, yes, there is something at least missing from your yaml most likely, and it's kind of important.  It's a data definition.  "Wha? Are those CD's in your backpack?  Don't try to install Oracle on my laptop!  Security!!"  No, no.  No need to worry, no relational databases here. 

The approach of kvse is to embed some lightweight "definition-like" syntax in your yaml to make it more expressive and extensible.  "I can do that already", you say.  Well, yes, but, this adds overhead, and we are not talking about strings vs. booleans, just to be clear: we are talking about your data and what it means to you, such as `environment`, `region`, or `collection`. 

Let's take a look at some techniques to add data definition into structured data. We know how relational databases do it, as a separate table basically, so let's not go over that, but for simple yaml files we might see a merely implicit data definition, or no definition at all - which is still implicit, because in the end we must know what the values represent in order to use them.

One could embed the data definition in various ways.

1. Fixed Column Definitions

In other words, like a relational database table.  This is *the* all time most favored approach. 

`ProjectName.Env.Region.Attribute`

`myproject.prod.us-west-2.vpc-id = vpc-1234`

Sure, that's an overly simplistic example, but, to put it to good use, it's going to get complicated fast.  What about things that don't have regions?  What about something that doesn't fit into an environment, like a role or data center?  We would either add this data into every environment, repeating ourselves in doing so, or have to break out a separate data set with a different structure, usually a different file also, to contain these different data.

This can work, but is extremely confining.  This may easily morph by adding dangling attributes (see #3 below) to overcome the limitations but doing so bloats the number of records while adding another level of complexity.

2. Adding Definition Keys within the Data

`myproject.environment.prod.region.us-west-2.vpcs.vpc-1`

This is similar in approach to LDAP and similarly verbose.  Unlike LDAP, here we have trouble keeping track of which elements are field selectors and which are the data.  The key lengths are bloated and we still need something that understands the schema.

3. Dangling Attributes

While not very elegant looking, this has the advantage that any element can be described in any number of ways.  A drawback would be in the level of effort in designing and maintaining an engine that understands how to work with this dynamic data definition approach.

`myproject.prod.is_env`

`myproject.prod.us-west-2.is_region`

`myproject.prod.us-west-2.vpc-1.is_vpc` 


### This Seems Hard - Is it Really Necessary?

Well, yes, assuming you will frequently need to know the answers to questions like:

- what regions does my project support?  
- what are all my environments?
- what is the path to the secret for the qa environment?
- what access lists are used in my project?

and so many more. 

Without either an implicit or explicit data definition we won't be able to get the answers to these questions - at least **not from our data** and **not programmatically** - two inescapable requirements of any automation at scale.

### How Does Kvse Add Data Definition While Keeping Things Simple?

1. Even before adding schema elements, Kvse can work with yaml you may already have
  - It will flatten your kvse-enabled yaml to key/value pairs allowing simpler consumption via key lookups
2. Kvse's data definition is inline
  - This means it's where you need to see and work with it during data manipulation - constructing and augmenting your yaml data.
3. Kvse supports dynamic names and structures
  - Any number of structures (schemas) can be designed, changed or extended in your yaml data.
4. Schema elements enable kvse to answer more advanced queries on the elements,
  - For example `project.environments` can produce a list of environments - no specialized query language needed
5. Kvse doesn't bloat your key space, because it removes the inline schema elements.

## Using Kvse

First note, that any array read by Kvse, currently must be for the purpose of embedding schema into the data elements of the array.  We can call these schema-enhanced lists.

- 'name' is a reserved term for use in schema-enhanced lists as a key in each object of the list.  If 'name' appears at all, it must appear in every object in the list.  It indicates implicitly that this is to replace metadata.
 
Maps are still just maps, the keys provide a single node in a dot-separated key space.

Map keys which are the parent of a schema-enhanced list act as placeholder schema elements, we call metadata keys.  These parent keys, or metadata keys, are what the 'name' values mentioned above replace during Kvse's rendering of yaml to key/value pairs.

The value of the metadata keys will be replaced, in turn, by each 'name' field value taken from the schema-enhanced list below it.

Example:

In the below scenario, the envs list is schema-enhanced since every element of the list contains a 'name' as a key field.  These names will be populated in place of the 'envs' element.  In essence, every 'name' field is naming a unique member of the 'envs' type.  

```
myproj:
  envs:
    - name: dev
      region: us-west-2
      account: 1234
    - name: test
      region: us-east-1
      account: 2345
```

kvse output:

```
myproj.dev.name=dev
myproj.dev.region=us-west-2
myproj.dev.account=1234
myproj.test.name=test
myproj.test.region=us-east-1
myproj.test.account=2345
```

We can repeat this pattern even inside a given schema-enhanced list.  In this next example,
we see the 'envs' schema-enhanced list as before, except we have added attributes in a schema-enhanced list under a newly defined meta-element, regions.

```
myproj:
  envs:
    - name: dev
      regions:
        - name: us-west-2
          vpcid: vpc-1234
          active: true
        - name: us-east-1
          vpcid:
          active: false
      account: 1234
    - name: test
      region: us-east-1
      account: 2345
```

kvse output:

```
myproj.dev.name=dev
myproj.dev.us-west-2.active=true
myproj.dev.us-west-2.name=us-west-2
myproj.dev.us-west-2.vpcid=vpc-1234
myproj.dev.us-east-1.name=us-east-1
myproj.dev.us-east-1.vpcid=null
myproj.dev.us-east-1.active=false
myproj.dev.account=1234
myproj.test.account=2345
myproj.test.name=test
myproj.test.region=us-east-1
```

## Build and Run

`make build`

`bin/kvse read myproj.envs `
